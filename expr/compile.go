package main

import (
	"fmt"

	"github.com/antonmedv/expr"
)

func main() {

	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf, // You can pass any functions.
	}
	code := `sprintf(greet, names[0])`

	// Compile code into bytecode. This step may be done
	// once and program may be reused.
	// Specify environment for type check.
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
