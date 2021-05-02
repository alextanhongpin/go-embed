// References: https://github.com/antonmedv/expr/blob/master/docs/Visitor-and-Patch.md
// This package demonstrates `Patch`. Implemented visitor
// can be applied before compileing AST to bytecode in
// expr.Compile function.
package main

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
)

func main() {
	env := map[string]interface{}{
		"list": []int{1, 2, 3},
	}
	code := `list[-1]` // will output 3
	program, err := expr.Compile(code, expr.Env(env), expr.Patch(&patcher{}))
	if err != nil {
		panic(err)
	}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}

type patcher struct{}

func (p *patcher) Enter(_ *ast.Node) {}
func (p *patcher) Exit(node *ast.Node) {
	n, ok := (*node).(*ast.IndexNode)
	if !ok {
		return
	}
	unary, ok := n.Index.(*ast.UnaryNode)
	if !ok {
		return
	}
	if unary.Operator == "-" {
		ast.Patch(&n.Index, &ast.BinaryNode{
			Operator: "-",
			Left:     &ast.BuiltinNode{Name: "len", Arguments: []ast.Node{n.Node}},
			Right:    unary.Node,
		})
	}
}
