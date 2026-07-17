# JSVM

基于 [goja](https://github.com/dop251/goja) 的 JavaScript 虚拟机封装，为 Go 嵌入 JavaScript 脚本执行能力，并提供一套 ESM
风格的模块化系统，将 Go 标准库以 JavaScript 接口的形式暴露给脚本代码。

## 特性

- **ESM 模块系统** — 脚本使用 `import` / `export` 语法，内部自动通过 [esbuild](https://github.com/evanw/esbuild) 转译并编译为
  goja 程序
- **Go 标准库导出** — 将 `net/http`、`strconv`、`time` 等标准库封装为 JS 模块，脚本中可直接使用
- **资源生命周期管理** — VM 关闭时自动按注册逆序清理所有资源（HTTP 服务、文件句柄等）
- **Context 联动取消** — VM 与 `context.Context` 绑定，父级取消时 VM 自动退出
- **JSON Tag 映射** — Go struct 字段通过 `json` tag 自动映射为 JS 的驼峰命名属性，无需手工转换

## 安装

```bash
go get github.com/xmx/jsvm
```

## 快速开始

```go
package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/xmx/jsvm"
	"github.com/xmx/jsvm/module/jsconsole"
	"github.com/xmx/jsvm/module/jshttp"
)

func main() {
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
```

## 内置模块

JSVM 内置了对以下 Go 标准库的封装，脚本中通过模块名 `import` 即可使用。

### 标准库封装

| 模块                                      | import 名  | 说明                            |
|-----------------------------------------|-----------|-------------------------------|
| [console](vitepress/modules/console.md) | `console` | 分级日志输出                        |
| [context](vitepress/modules/context.md) | `context` | 上下文管理：超时、取消、值传递               |
| [io](vitepress/modules/io.md)           | `io`      | 读写操作：copy、readAll、limitReader |
| [os](vitepress/modules/os.md)           | `os`      | 操作系统交互：环境变量、进程信息              |
| [runtime](vitepress/modules/runtime.md) | `runtime` | 运行时信息：CPU、goroutine、内存统计      |
| [strconv](vitepress/modules/strconv.md) | `strconv` | 字符串与数字互转                      |
| [strings](vitepress/modules/strings.md) | `strings` | 字符串操作：查找、替换、分割                |
| [time](vitepress/modules/time.md)       | `time`    | 时间处理：Duration、Time、Location   |

### 网络相关

| 模块                                                     | import 名            | 说明                                  |
|--------------------------------------------------------|---------------------|-------------------------------------|
| [net/http](vitepress/modules/net/http.md)              | `net/http`          | HTTP 客户端（fetch）与服务端（listenAndServe） |
| [net/http/httputil](vitepress/modules/net/httputil.md) | `net/http/httputil` | HTTP 工具：反向代理                        |
| [net/url](vitepress/modules/net/url.md)                | `net/url`           | URL 解析与查询参数构建                       |

## 项目结构

```
jsvm/
├── vm.go             # VM 核心：创建、执行、生命周期
├── module.go         # ModuleExporter / ModuleExports 接口定义
├── compile.go        # esbuild 编译（JSX/TS → CommonJS → goja.Program）
├── mapper.go         # Go struct 字段的 tagName 映射器
├── cleanup.go        # 资源清理管理器
├── utils.go          # 工具函数（IsNullish 等）
├── module/           # 内置 Go 标准库封装模块
│   ├── jsconsole/
│   ├── jscontext/
│   ├── jshttp/
│   ├── jshttputil/
│   ├── jsio/
│   ├── jsos/
│   ├── jsruntime/
│   ├── jsstrconv/
│   ├── jsstrings/
│   ├── jstime/
│   └── jsurl/
├── examples/         # 使用示例
├── docs/             # VitePress 文档站
└── vitepress/        # 文档源文件
```

## 自定义模块

实现 `ModuleExporter` 接口即可扩展 JSVM：

```go
package mymodule

import (
	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

type myModule struct{}

func New() jsvm.ModuleExporter { return &myModule{} }

func (m *myModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	return jsvm.ModuleExports{
		Name: "mymodule",
		Default: map[string]any{
			"hello": func(name string) string {
				return "Hello, " + name + "!"
			},
		},
	}
}
```

脚本中使用：

```javascript
import {hello} from 'mymodule'

console.log(hello('JSVM'))  // "Hello, JSVM!"
```

## 资源管理

VM 中开启的服务可通过 `AddCleaner` 注册，VM 关闭时自动清理：

```go
srv := &http.Server{Addr: ":8080"}
cln, ok := vm.AddCleaner(srv)
// ok = true: 注册成功，VM 关闭时自动 srv.Close()
// ok = false: VM 已关闭，cln.Close() 仍可调用
```
