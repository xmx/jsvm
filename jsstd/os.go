package jsstd

import (
	"os"

	"github.com/xmx/jsvm"
)

type osModule struct {
	vm *jsvm.VM
}

func NewOS() jsvm.ModuleLoader {
	return &osModule{}
}

func (m *osModule) LoadModule(vm *jsvm.VM, opts jsvm.LoadModuleOptions) (string, map[string]any, error) {
	m.vm = vm
	vals := map[string]any{
		"args":          os.Args,
		"executable":    os.Executable,
		"environ":       os.Environ,
		"hostname":      os.Hostname,
		"userConfigDir": os.UserConfigDir,
		"userCacheDir":  os.UserCacheDir,
		"userHomeDir":   os.UserHomeDir,
		"getenv":        os.Getenv,
		"geteuid":       os.Geteuid,
		"getgid":        os.Getgid,
		"getppid":       os.Getppid,
		"getpid":        os.Getpid,
		"getwd":         os.Getwd,
	}
	if opts.Unsafe {
		vals["unsetenv"] = os.Unsetenv
	}

	return "os", vals, nil
}
