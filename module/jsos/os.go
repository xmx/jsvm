package jsos

import (
	"os"

	"github.com/xmx/jsvm"
)

// osModule 是 os 模块的实现，将 Go 标准库 os 包的常用函数导出到 JS。
type osModule struct {
	vm *jsvm.VM
}

// New 创建 os 模块实例。
func New() jsvm.ModuleExporter {
	return new(osModule)
}

// ModuleExports 注册 os 模块，导出：
//   - args：命令行参数列表
//   - executable：当前可执行文件路径
//   - environ：所有环境变量
//   - hostname：主机名
//   - userConfigDir/userCacheDir/userHomeDir：用户目录
//   - getenv：读取环境变量
//   - geteuid/getgid/getppid/getpid：进程信息
//   - getwd：当前工作目录
func (m *osModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm
	vals := map[string]any{
		"args":          os.Args,          // 命令行参数，args[0] 为程序名
		"executable":    os.Executable,    // 返回当前可执行文件路径
		"environ":       os.Environ,       // 返回所有环境变量（格式："key=value"）
		"hostname":      os.Hostname,      // 返回主机名
		"userConfigDir": os.UserConfigDir, // 用户配置目录（如 ~/.config）
		"userCacheDir":  os.UserCacheDir,  // 用户缓存目录（如 ~/.cache）
		"userHomeDir":   os.UserHomeDir,   // 用户主目录
		"getenv":        os.Getenv,        // 获取指定环境变量的值
		"geteuid":       os.Geteuid,       // 当前有效用户 ID
		"getgid":        os.Getgid,        // 当前组 ID
		"getppid":       os.Getppid,       // 父进程 PID
		"getpid":        os.Getpid,        // 当前进程 PID
		"getwd":         os.Getwd,         // 当前工作目录
	}

	return jsvm.ModuleExports{
		Name:    "os",
		Default: vals,
	}
}
