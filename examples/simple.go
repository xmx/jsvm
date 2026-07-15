package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/xmx/jsvm"
	"github.com/xmx/jsvm/jsstd"
)

func main() {
	ctx := context.Background()
	log := slog.Default()

	vm := jsvm.NewVM(ctx, log)
	vm.RegisterModules([]jsvm.ModuleExporter{jsstd.NewConsole(), jsstd.NewHTTP()})

	time.AfterFunc(10*time.Second, func() {
		fmt.Println("=========================")
		vm.Cancel()
	})

	val, err := vm.RunScript("hi", `
import {A} from 'console'
import http from 'net/http'

const mux = http.newServeMux()
mux.handleFunc('/ping', (w, r) => {
	A.log('收到请求 ' + r.remoteAddr)
	w.write('PONG')
})

http.listenAndServe(':8088', mux)
`)
	fmt.Println(err)
	fmt.Println(val)
}
