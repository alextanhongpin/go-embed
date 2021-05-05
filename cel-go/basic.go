// This can be used as the basics for rule engine.
package main

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
)

func main() {
	{
		expr := `data.name == "john"`
		in := map[string]interface{}{
			"name": "john",
		}
		out, err := eval(expr, in)
		if err != nil {
			panic(err)
		}
		fmt.Println(out)
	}

	{
		expr := `data.names.all(n, n.startsWith('a')) && size(data.names) == 2`
		in := map[string]interface{}{
			"names": []string{"aloha", "alpha"},
		}
		out, err := eval(expr, in)
		if err != nil {
			panic(err)
		}
		fmt.Println(out)
	}
}

func eval(expr string, in map[string]interface{}) (ref.Val, error) {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("data", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)

	ast, iss := env.Compile(expr)
	if iss != nil && iss.Err() != nil {
		return nil, iss.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"data": in,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
