package jsvm

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
)

// VM 是基于 goja 的 JavaScript 虚拟机。
// 提供生命周期管理、模块加载和资源清理功能。
type VM struct {
	log     *slog.Logger             // 日志记录器
	cleaner cleanManager             // 资源清理管理器
	rt      *goja.Runtime            // goja 运行时
	modules map[string]goja.Value    // 已加载的模块缓存（模块名 -> goja 值）
	exports map[string]ModuleExports // 已注册的模块导出（模块名 -> 导出定义）
	used    atomic.Bool              // VM 是否已执行过脚本（单次执行保护）
	oncec   func()                   // 确保 close 只执行一次
	ctx     context.Context          // VM 关联的上下文
	cancel  context.CancelFunc       // 取消函数，调用后 VM 停止运行
}

// NewVM 创建新的虚拟机实例，包含运行时、模块系统和资源清理机制。
// 当 parent 被取消或调用 Cancel 时，所有注册的资源将按注册顺序逆序关闭。
func NewVM(parent context.Context, log *slog.Logger) *VM {
	rt := goja.New()
	rt.SetFieldNameMapper(tagMapper("json"))

	ctx, cancel := context.WithCancel(parent)

	vm := &VM{
		log:     log,
		rt:      rt,
		cleaner: new(cleanMapManager),
		modules: make(map[string]goja.Value, 16),
		exports: make(map[string]ModuleExports, 16),
		ctx:     ctx,
		cancel:  cancel,
	}
	vm.oncec = sync.OnceFunc(vm.close)
	context.AfterFunc(ctx, vm.oncec)  // 上下文取消时自动执行清理
	_ = rt.Set("require", vm.require) // 注入 require 函数供脚本使用

	return vm
}

// Logger 返回 VM 的日志记录器。
func (vm *VM) Logger() *slog.Logger { return vm.log }

// Context 返回 VM 的上下文。VM 关闭时此上下文会被取消。
func (vm *VM) Context() context.Context { return vm.ctx }

// RegisterModules 批量注册模块到 VM。
// 每个模块必须实现 ModuleExporter 接口。
func (vm *VM) RegisterModules(mods []ModuleExporter) {
	for _, mod := range mods {
		if mod != nil {
			exp := mod.ModuleExports(vm)
			vm.exports[exp.Name] = exp
		}
	}
}

// AddCleaner 注册资源，VM 关闭时自动调用其 Close 方法。
// 返回句柄和是否成功标志：
//   - 成功(true)：调用句柄的 Close() 会关闭资源并从清理列表移除；
//   - 失败(false)：VM 已关闭，句柄的 Close() 仍可直接关闭资源，但 Unregister() 返回 nil。
func (vm *VM) AddCleaner(c io.Closer) (CleanHandle, bool) {
	ch := &fallbackCleanHandle{cm: vm.cleaner}
	if id, succ := vm.cleaner.register(c); succ {
		ch.id = id
		return ch, true
	}

	ch.fb = c // 降级：直接关闭

	return ch, false
}

// RunProgram 执行已编译的程序，只能调用一次。
// 若 VM 已执行过或已关闭，返回错误。
// 执行完成后 VM 自动取消，不可再次执行。
// 使用 Compile 创建程序。
func (vm *VM) RunProgram(pgm *goja.Program) (goja.Value, error) {
	return vm.rt.RunProgram(pgm)
}

// RunScript 编译并执行 JavaScript 代码，只能调用一次。
// name 参数用于错误提示和源码映射。
// 执行完成后 VM 自动取消。
func (vm *VM) RunScript(name, code string) (goja.Value, error) {
	pgm, err := Compile(name, code)
	if err != nil {
		return nil, err
	}

	return vm.RunProgram(pgm)
}

// Throw 在 JavaScript 运行时抛出 Go 错误。
// 若 err 已是 goja Exception，直接 panic 该异常；
// 否则包装为 goja GoError 并 panic（保留调用栈）。
func (vm *VM) Throw(err error) {
	//goland:noinspection GoTypeAssertionOnErrors
	if e, ok := err.(*goja.Exception); ok {
		panic(e)
	}
	panic(vm.rt.NewGoError(err)) // 保留调用栈，优于 rt.ToValue
}

// Cancel 关闭 VM：取消上下文、中断正在运行的脚本、关闭所有注册资源。
// 可多次调用，并发调用会阻塞直到首次关闭完成。
func (vm *VM) Cancel() {
	vm.cancel()
	vm.oncec()
}

// close 执行实际清理：中断运行时并关闭所有资源。
func (vm *VM) close() {
	vm.rt.Interrupt(context.Canceled)
	vm.cleaner.closeAll()
}

// require 是注入到运行时中的 require 函数，用于按模块名加载已注册模块。
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

// resolve 按名称查找模块：先从缓存读取，再从 exports 中首次解析并缓存。
// 返回值不为 nil 表示已找到；err 不为 nil 表示解析过程出错。
func (vm *VM) resolve(name string) (goja.Value, error) {
	if val, ok := vm.modules[name]; ok {
		return val, nil
	}

	exp, ok := vm.exports[name]
	if !ok {
		return nil, nil // 模块不存在
	}

	esm := exp.toESM()
	val := vm.rt.ToValue(esm)
	vm.modules[name] = val // 缓存到 modules，避免重复转换

	return val, nil
}
