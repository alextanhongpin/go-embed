# go-embed


Why embed another language in golang?
- WIP

What other language do we have?

- expr
- python
- two lib for lua, [gopher-lua](https://github.com/yuin/gopher-lua) and [go-lua](https://github.com/Shopify/go-lua)
- [otto](https://github.com/robertkrimen/otto) allows ES5 JavaScript only, see v8g0
- [v8go](https://github.com/rogchap/v8go), probably can be used to emulate Cloudflare Workers (which are most likely using containers)
- [cuego](https://pkg.go.dev/cuelang.org/go@v0.3.2/cuego) as well as [cuelang](https://cuelang.org/docs/references/)
- [tengo](https://github.com/d5/tengo)
- [go+](https://github.com/goplus/gop) a language for engineering, STEM education, and data science
- [jsonnet](https://jsonnet.org/learning/tutorial.html)
- [starlark](https://github.com/bazelbuild/starlark), and [starlark-go](https://github.com/google/starlark-go)
- [cel-go](https://github.com/google/cel-go)
- [go-bexpr](https://github.com/hashicorp/go-bexpr) from Hashicorp
- [jsonata](https://jsonata.org/)
- https://www.kcl-lang.io/docs/user_docs/getting-started/kcl-quick-start

For javascript, there's vm and [vm2](https://github.com/patriksimek/vm2), and [isolated-vm](https://github.com/laverdet/isolated-vm) that is used by [temporal](https://docs.temporal.io/blog/intro-to-isolated-vm/).


Another language that is similar to go is [vlang](https://vlang.io/compare#go).

https://dhall-lang.org

## References

1. https://otm.github.io/2015/07/embedding-lua-in-go/
2. https://go.libhunt.com/categories/485-embeddable-scripting-languages


## go-jsonnet

An alternative to modify json dynamically in golang since jsonata port doesn't exist.

```go
package main

import (
	"fmt"
	"log"

	"github.com/google/go-jsonnet"
)

func main() {
	vm := jsonnet.MakeVM()

	snippet := `{
		person1: {
		    name: "Alice",
		    welcome: "Hello " + std.extVar("name").name + "!",
		},
		person2: self.person1 { name: "Bob" },
	}`

	// vm.ExtVar("sth", "else") // This is only string
	vm.ExtCode("name", `{"name": "john"}`) // Object is allowed
	jsonStr, err := vm.EvaluateAnonymousSnippet("example1.jsonnet", snippet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jsonStr)
	/*
	   {
	     "person1": {
	         "name": "Alice",
	         "welcome": "Hello Alice!"
	     },
	     "person2": {
	         "name": "Bob",
	         "welcome": "Hello Bob!"
	     }
	   }
	*/
}
```
