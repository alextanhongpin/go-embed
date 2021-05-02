package main

import (
	"fmt"

	"github.com/antonmedv/expr"
)

func main() {
	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}
	// Evaluate expressions, as well as type checking such
	// expression.
	out, err := expr.Eval("foo + bar", env)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
