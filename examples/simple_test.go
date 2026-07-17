package examples_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/xmx/jsvm"
	"github.com/xmx/jsvm/module/jsconsole"
	"github.com/xmx/jsvm/module/jscontext"
	"github.com/xmx/jsvm/module/jshttp"
	"github.com/xmx/jsvm/module/jshttputil"
	"github.com/xmx/jsvm/module/jsio"
	"github.com/xmx/jsvm/module/jsos"
	"github.com/xmx/jsvm/module/jsruntime"
	"github.com/xmx/jsvm/module/jsstrconv"
	"github.com/xmx/jsvm/module/jsstrings"
	"github.com/xmx/jsvm/module/jstime"
	"github.com/xmx/jsvm/module/jsurl"
)

func TestSimple(t *testing.T) {
	const filename = "time.js"
	raw, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	code := string(raw)

	ctx := context.Background()
	log := slog.Default()
	mods := []jsvm.ModuleExporter{
		jsconsole.New(),
		jscontext.New(),
		jshttp.New(),
		jshttputil.New(),
		jsio.New(),
		jsos.New(),
		jsruntime.New(),
		jsstrconv.New(),
		jsstrings.New(),
		jstime.New(),
		jsurl.New(),
	}

	vm := jsvm.NewVM(ctx, log)
	vm.RegisterModules(mods)
	time.AfterFunc(60*time.Second, vm.Cancel)

	_, err = vm.RunScript(filename, code)
	t.Log(err)
}

func TestHTTP(t *testing.T) {
	vm := jsvm.NewVM(context.Background(), slog.Default())
	vm.RegisterModules([]jsvm.ModuleExporter{
		jsconsole.New(),
		jshttp.New(),
	})

	_, err := vm.RunScript("main.js", `
        import console from 'console'
        import http from 'net/http'

        const res = http.fetch('https://httpbin.io/get')
        console.log(res.text())
    `)
	if err != nil {
		fmt.Println("error:", err)
	}
}
