// References: https://github.com/antonmedv/expr/blob/master/docs/Optimizations.md
package main

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

func main() {

	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	program, err := expr.Compile("foo + bar", expr.Env(env))
	if err != nil {
		panic(err)
	}

	// Reuse this vm instance between runs.
	v := vm.VM{}
	out, err := v.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
