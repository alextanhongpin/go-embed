// References: https://github.com/antonmedv/expr/blob/master/docs/Visitor-and-Patch.md
package main

import (
	"fmt"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/parser"
)

type visitor struct {
	identifiers []string
}

func (v *visitor) Enter(node *ast.Node) {}
func (v *visitor) Exit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok {
		v.identifiers = append(v.identifiers, n.Value)
	}
}

func main() {
	tree, err := parser.Parse("foo + bar")
	if err != nil {
		panic(err)
	}

	visitor := &visitor{}
	ast.Walk(&tree.Node, visitor)
	fmt.Printf("%v", visitor.identifiers) // outputs [foo bar]
}
