package examples_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/xmx/jsvm"
	"github.com/xmx/jsvm/jsstd"
)

func TestSimple(t *testing.T) {
	raw, err := os.ReadFile("main.js")
	if err != nil {
		t.Fatal(err)
	}
	code := string(raw)

	ctx := context.Background()
	log := slog.Default()
	mods := []jsvm.ModuleExporter{
		jsstd.NewConsole(),
		jsstd.NewHTTP(),
	}

	vm := jsvm.NewVM(ctx, log)
	vm.RegisterModules(mods)
	time.AfterFunc(60*time.Second, vm.Cancel)

	_, err = vm.RunScript("main.js", code)
	t.Log(err)
}
