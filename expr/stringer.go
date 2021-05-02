// References: https://github.com/antonmedv/expr/blob/master/docs/Visitor-and-Patch.md
package main

import (
	"fmt"
	"reflect"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
)

func main() {
	code := `Price == "$100"`

	program, err := expr.Compile(code, expr.Env(Env{}), expr.Patch(&stringerPatcher{}))
	if err != nil {
		panic(err)
	}

	env := Env{100_00}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}

type Env struct {
	Price Price
}

type Price int

func (p Price) String() string {
	return fmt.Sprintf("$%v", int(p)/100)
}

var stringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

type stringerPatcher struct{}

func (*stringerPatcher) Enter(_ *ast.Node) {}
func (*stringerPatcher) Exit(node *ast.Node) {
	t := (*node).Type()
	if t == nil {
		return
	}
	if t.Implements(stringer) {
		ast.Patch(node, &ast.MethodNode{
			Node:   *node,
			Method: "String",
		})
	}
}
