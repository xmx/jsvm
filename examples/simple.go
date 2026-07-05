package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/xmx/jsvm"
	"github.com/xmx/jsvm/jsstd"
)

func main() {
	mods := []jsvm.ModuleLoader{
		jsstd.NewOS(),
		jsstd.NewConsole(),
	}

	log := slog.Default()
	ctx := context.Background()
	vm := jsvm.NewVM(ctx, log)
	_ = vm.AddModules(mods, jsvm.LoadModuleOptions{Unsafe: true})

	val, err := vm.RunScript("test.js", script)
	fmt.Println(err)
	fmt.Println(val)
}

const script = `
import os from 'os'
import console from 'console'

const pid = os.getpid()
console.log(pid, 'INF')

const fn = ()=>{}
console.log(fn)
`
