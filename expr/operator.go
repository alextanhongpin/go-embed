// This example demonstrates how to do operator overriding with expr.
package main

import (
	"fmt"
	"time"

	"github.com/antonmedv/expr"
)

func main() {
	code := `(Now() - CreatedAt).Hours() / 24 / 365`

	// We can define options before compiling.
	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator("-", "Sub"), // Override `-` with function `Sub`.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		panic(err)
	}
	env := Env{
		CreatedAt: time.Date(1987, time.November, 24, 20, 0, 0, 0, time.UTC),
	}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}

type Env struct {
	datetime
	CreatedAt time.Time
}

// Functions may be defined on embedded structs as well.
type datetime struct{}

func (datetime) Now() time.Time                   { return time.Now() }
func (datetime) Sub(a, b time.Time) time.Duration { return a.Sub(b) }
