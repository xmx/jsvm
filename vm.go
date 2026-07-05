package jsvm

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync/atomic"

	"github.com/dop251/goja"
)

// VM is a JavaScript virtual machine backed by goja.
// It provides lifecycle management, module loading, and cleanup.
type VM struct {
	log     *slog.Logger
	cleaner cleaner
	rt      *goja.Runtime
	modules map[string]goja.Value
	closed  atomic.Bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewVM creates a new VM. When parent is cancelled or Cancel is called,
// all registered CleanHandles are closed in reverse registration order.
func NewVM(parent context.Context, log *slog.Logger) *VM {
	rt := goja.New()
	mapper := goja.TagFieldNameMapper("json", true)
	rt.SetFieldNameMapper(mapper)
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
func (vm *VM) AddModules(mods []ModuleLoader, opts LoadModuleOptions) error {
	rt := vm.rt
	for _, lm := range mods {
		pkg, vals, err := lm.LoadModule(vm, opts)
		if err != nil {
			return err
		}
		vm.modules[pkg] = rt.ToValue(vals)
	}

	return nil
}

// AddCleaner registers a resource to be closed when the VM shuts down.
// Returns a CleanHandle and true on success, or nil and false if the VM is closed.
func (vm *VM) AddCleaner(c io.Closer) (CleanHandle, bool) {
	if vm.closed.Load() {
		return nil, false
	}
	id := vm.cleaner.register(c)

	return &cleanHandle{
		id:  id,
		cln: vm.cleaner,
	}, true
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
	if val, ok := vm.modules[name]; ok {
		return val
	}

	err := errors.New("module not found: " + name)
	panic(vm.rt.NewGoError(err))
}
