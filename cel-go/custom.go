package main

import (
	"log"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter/functions"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type Test struct {
	a, b int
}

func (t *Test) Add() int {
	return t.a + t.b
}

func (t *Test) Subtract() int {
	return t.a - t.b
}

// The CEL type to represent Test.
var TestType = types.NewTypeValue("Test", traits.ReceiverType)

func (t Test) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	panic("not required")
}

func (t Test) ConvertToType(typeVal ref.Type) ref.Val {
	panic("not required")
}

func (t Test) Equal(other ref.Val) ref.Val {
	o, ok := other.Value().(Test)
	if ok {
		return types.Bool(o == t)
	}
	return types.ValOrErr(other, "%v is not of type Test", other)
}

func (t Test) Type() ref.Type {
	return TestType
}

func (t Test) Value() interface{} {
	return t
}

func (t Test) Receive(function string, overload string, args []ref.Val) ref.Val {
	if function == "Add" {
		return types.Int(t.Add())
	} else if function == "Subtract" {
		return types.Int(t.Subtract())
	}
	return types.ValOrErr(TestType, "no such function - %s", function)
}

func (t *Test) HasTrait(trait int) bool {
	return trait == traits.ReceiverType
}

func (t *Test) TypeName() string {
	return TestType.TypeName()
}

type customTypeAdapter struct{}

func (c customTypeAdapter) NativeToValue(value interface{}) ref.Val {
	val, ok := value.(Test)
	if ok {
		return val
	}
	// Let the default adapter handle other cases.
	return types.DefaultTypeAdapter.NativeToValue(value)
}

func main() {
	env, err := cel.NewEnv(
		cel.CustomTypeAdapter(&customTypeAdapter{}),
		cel.Declarations(
			decls.NewVar("test", decls.NewObjectType("Test")),
			decls.NewFunction("MulBy3", decls.NewOverload("mulby3_int", []*expr.Type{decls.Int}, decls.Int)),
			decls.NewFunction("Add", decls.NewInstanceOverload("test_add", []*expr.Type{decls.NewObjectType("Test")}, decls.Int)),
			decls.NewFunction("Subtract", decls.NewInstanceOverload("test_subtract", []*expr.Type{decls.NewObjectType("Test")}, decls.Int)),
		),
	)
	if err != nil {
		panic(err)
	}

	parsed, issues := env.Parse(`test.Add() == 3 && test.Subtract() == -1 && MulBy3(9) == 27`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("parse error: %v", issues.Err())
	}

	checked, issues := env.Check(parsed)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	globalFunctions := cel.Functions(
		&functions.Overload{
			Operator: "MulBy3",
			Unary: func(lhs ref.Val) ref.Val {
				return types.Int(3 * lhs.Value().(int64))
			},
		},
	)

	prg, err := env.Program(checked, globalFunctions)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"test": Test{a: 1, b: 2},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got output: %v", out)
}
