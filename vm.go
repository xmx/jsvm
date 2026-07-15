package jsvm

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
)

// ErrExecuted is returned when RunProgram or RunScript is called more than once.
var ErrExecuted = errors.New("vm: already executed")

// VM is a JavaScript virtual machine backed by goja.
// It provides lifecycle management, module loading, and cleanup.
type VM struct {
	log       *slog.Logger
	cleaner   cleanManager
	rt        *goja.Runtime
	exports   map[string]ModuleExports
	used      atomic.Bool
	onceClose func()
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewVM creates a new VM with a runtime, module system, and cleanup mechanism.
// When parent is cancelled or Cancel is called, all registered resources
// are automatically closed in reverse registration order.
func NewVM(parent context.Context, log *slog.Logger) *VM {
	rt := goja.New()
	rt.SetFieldNameMapper(tagMapper("json"))

	ctx, cancel := context.WithCancel(parent)

	vm := &VM{
		log:     log,
		rt:      rt,
		cleaner: new(cleanMapManager),
		exports: make(map[string]ModuleExports, 16),
		ctx:     ctx,
		cancel:  cancel,
	}
	vm.onceClose = sync.OnceFunc(vm.close)
	context.AfterFunc(ctx, vm.onceClose)
	_ = rt.Set("require", vm.require)

	return vm
}

// Logger returns the VM's logger.
func (vm *VM) Logger() *slog.Logger { return vm.log }

// Context returns the VM's context. It is cancelled when the VM is closed.
func (vm *VM) Context() context.Context { return vm.ctx }

func (vm *VM) RegisterModules(mods []ModuleExporter) {
	for _, mod := range mods {
		if mod == nil {
			continue
		}
		exp := mod.ModuleExports(vm)
		vm.exports[exp.Name] = exp
	}
}

// AddCleaner registers a resource to be closed when the VM is closed.
// Returns a handle and a success flag. On success (true), calling handle.Close()
// closes the resource and removes it from cleanup; handle.Unregister() returns
// the resource. On failure (false, e.g., VM already closed), handle.Close()
// still closes the resource directly, but Unregister() returns nil.
func (vm *VM) AddCleaner(c io.Closer) (CleanHandle, bool) {
	ch := &fallbackCleanHandle{cm: vm.cleaner}
	if id, succ := vm.cleaner.register(c); succ {
		ch.id = id
		return ch, true
	}

	ch.fb = c

	return ch, false
}

// RunProgram executes a compiled program once. Returns ErrExecuted if the VM has
// already been used, or context.Canceled if the VM is closed.
// After execution completes, the VM is automatically cancelled.
// Use Compile to create programs from source code.
func (vm *VM) RunProgram(pgm *goja.Program) (goja.Value, error) {
	if !vm.used.CompareAndSwap(false, true) {
		return nil, ErrExecuted
	}
	defer vm.Cancel()

	return vm.rt.RunProgram(pgm)
}

// RunScript compiles and executes JavaScript code once. The name parameter is used
// for error messages and source maps. After execution, the VM is automatically cancelled.
func (vm *VM) RunScript(name, code string) (goja.Value, error) {
	pgm, err := Compile(name, code)
	if err != nil {
		return nil, err
	}

	return vm.RunProgram(pgm)
}

func (vm *VM) Throw(err error) {
	//goland:noinspection GoTypeAssertionOnErrors
	if e, ok := err.(*goja.Exception); ok {
		panic(e)
	}
	panic(vm.rt.NewGoError(err)) // this catches the stack unlike rt.ToValue
}

// Cancel shuts down the VM. It cancels the VM's context, interrupts any running
// script, and closes all registered resources. Safe to call multiple times;
// concurrent and subsequent calls block and all receive the same error from the
// first shutdown.
func (vm *VM) Cancel() {
	vm.cancel()
	vm.onceClose()
}

func (vm *VM) close() {
	vm.rt.Interrupt(context.Canceled)
	vm.cleaner.closeAll()
}

func (vm *VM) require(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	obj, err := vm.resolve(name)
	if obj != nil {
		return obj
	}
	if err != nil {
		vm.Throw(err)
	}

	panic(vm.rt.NewTypeError("cannot find module '%s'.", name))
}

func (vm *VM) resolve(name string) (goja.Value, error) {
	exp, ok := vm.exports[name]
	if !ok {
		return nil, nil
	}

	esm := exp.toESM()
	ret := vm.rt.ToValue(esm)

	return ret, nil
}
