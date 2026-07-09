package jsstd

import (
	"os"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

func NewOS() jsvm.Module {
	return &osModule{}
}

type osModule struct {
	vm *jsvm.VM
}

func (m *osModule) Name() string {
	return "os"
}

func (m *osModule) Load(vm *jsvm.VM, exports *goja.Object) error {
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

	return jsvm.SetExports(exports, vals)
}
