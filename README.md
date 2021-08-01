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
- [jsonnet](https://jsonnet.org/learning/tutorial.html)
- [starlark](https://github.com/bazelbuild/starlark), and [starlark-go](https://github.com/google/starlark-go)
- [cel-go](https://github.com/google/cel-go)
- [go-bexpr](https://github.com/hashicorp/go-bexpr) from Hashicorp

For javascript, there's vm and [vm2](https://github.com/patriksimek/vm2)

https://dhall-lang.org

## References

1. https://otm.github.io/2015/07/embedding-lua-in-go/
2. https://go.libhunt.com/categories/485-embeddable-scripting-languages
