package jsvm

import (
	"context"
	"io"
	"log/slog"
	"sync/atomic"

	"github.com/dop251/goja"
)

// VM is a JavaScript virtual machine backed by goja.
// It provides lifecycle management, module loading, and cleanup.
type VM struct {
	log      *slog.Logger
	cleaner  cleaner
	rt       *goja.Runtime
	modules  map[string]goja.Value
	modloads map[string]Module
	closed   atomic.Bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewVM creates a new VM. When parent is cancelled or Cancel is called,
// all registered CleanHandles are closed in reverse registration order.
func NewVM(parent context.Context, log *slog.Logger) *VM {
	rt := goja.New()
	rt.SetFieldNameMapper(tagMapper("json"))
	cln := newCleanerMap(log)
	modules := make(map[string]goja.Value, 16)
	ctx, cancel := context.WithCancel(parent)

	vm := &VM{
		log:     log,
		rt:      rt,
		cleaner: cln,
		modules: modules,
		ctx:     ctx,
		cancel:  cancel,
	}
	_ = rt.Set("require", vm.require)
	context.AfterFunc(ctx, vm.closeNopError)

	return vm
}

// Logger returns the VM's logger.
func (vm *VM) Logger() *slog.Logger {
	return vm.log
}

// Context returns the VM's context. It is cancelled when the VM is closed.
func (vm *VM) Context() context.Context {
	return vm.ctx
}

// AddModules registers modules from the given loaders. The returned
// values are accessible to JS code via require(name).
func (vm *VM) AddModules(mods []Module) {
	if vm.modloads == nil {
		vm.modloads = make(map[string]Module, len(mods))
	}
	for _, mod := range mods {
		name := mod.Name()
		vm.modloads[name] = mod
	}
}

// AddCleaner registers a resource to be closed when the VM shuts down.
// Returns a CleanHandle and true on success, or nil and false if the VM is closed.
func (vm *VM) AddCleaner(c io.Closer) (CleanHandle, bool) {
	cln := &cleanHandle{cln: vm.cleaner}
	if vm.closed.Load() {
		cln.back = c
		return cln, false
	}
	cln.id = vm.cleaner.register(c)

	return cln, true
}

// RunProgram executes a compiled program. Returns context.Canceled if the VM is closed.
func (vm *VM) RunProgram(pgm *goja.Program) (goja.Value, error) {
	if vm.closed.Load() {
		return nil, context.Canceled
	}

	return vm.rt.RunProgram(pgm)
}

// RunScript compiles and runs JavaScript code. The name is used for error
// messages and source maps.
func (vm *VM) RunScript(name, code string) (goja.Value, error) {
	pgm, err := Compile(name, code)
	if err != nil {
		return nil, err
	}

	return vm.RunProgram(pgm)
}

// Cancel shuts down the VM. It interrupts any running script, closes
// all registered resources, and cancels the VM's context.
func (vm *VM) Cancel() error {
	vm.cancel()
	return vm.close()
}

func (vm *VM) closeNopError() {
	if err := vm.close(); err != nil {
		vm.log.Debug("vm closed error", "err", err)
	}
}

func (vm *VM) close() error {
	if !vm.closed.CompareAndSwap(false, true) {
		return context.Canceled
	}

	vm.rt.Interrupt(context.Canceled)
	vm.cleaner.execute()

	return nil
}

func (vm *VM) require(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	obj, err := vm.resolve(name)
	if obj != nil {
		return obj
	}
	if err != nil {
		if _, ok := err.(*goja.Exception); !ok {
			panic(vm.rt.NewGoError(err))
		}
		panic(err)
	}

	panic(vm.rt.NewTypeError("cannot find module '%s'.", name))
}

func (vm *VM) resolve(name string) (goja.Value, error) {
	val, ok := vm.modules[name]
	if ok {
		return val, nil
	}

	ld, yes := vm.modloads[name]
	if !yes {
		return nil, nil
	}

	// 忽略各种复杂的实现，只关心 exports
	exp := vm.rt.NewObject()
	if err := ld.Load(vm, exp); err != nil {
		return nil, err
	}

	if vm.modules == nil {
		vm.modules = make(map[string]goja.Value, 8)
	}
	vm.modules[name] = exp

	return exp, nil
}
