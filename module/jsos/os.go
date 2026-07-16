package jsos

import (
	"os"

	"github.com/xmx/jsvm"
)

type osModule struct {
	vm *jsvm.VM
}

func New() jsvm.ModuleExporter {
	return new(osModule)
}

func (m *osModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
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

	return jsvm.ModuleExports{
		Name:    "os",
		Default: vals,
	}
}
