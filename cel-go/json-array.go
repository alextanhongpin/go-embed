package main

import (
	"encoding/json"
	"log"

	//structpb "google.golang.org/protobuf/types/known/structpb"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type Payload struct {
	Strs []string          `json:"strs"`
	Data map[string]string `json:"data"`
}

type MyList struct {
	Items []MyStruct `json:"items"`
}

type MyStruct struct {
	Num     int64   `json:"num"`
	Str     string  `json:"str"`
	Payload Payload `json:"payload"`
}

// An example of using native JSON-friendly Go struct as a type and input in a Cel program.
func main() {
	// NOTE: myStruct.num == 10 will fail with error: no such overload.
	//https://github.com/google/cel-go/issues/203
	// Unique str
	//filter := `!data.items.exists(
	//item1,
	//data.items
	//.filter(
	//item,
	//int(item.num) != int(item1.num)
	//)
	//.exists(
	//item2,
	//string(item2.str) == string(item1.str)
	//)
	//)` // works.

	//filter := `uniq(data.items.map(item, item.str))` // works.
	//filter := `uniq(data.items.map(item, item.num))` // works.
	filter := `data.items.map(item, item.str).uniq()` // works.
	//filter := `uniq(data.items.map(item, item.payload))` // does not work.

	//filter := `data.items.all(item, item.str.startsWith("h"))` // works.

	myList := MyList{
		Items: []MyStruct{
			MyStruct{
				Num: 10,
				Str: "hello",
				Payload: Payload{
					Data: map[string]string{"world": "foobar"},
					Strs: []string{"banana"},
				},
			},
			MyStruct{
				Num: 11,
				Str: "hello",
				Payload: Payload{
					Data: map[string]string{"world": "foobar"},
					Strs: []string{"banana"},
				},
			},
		},
	}
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewFunction("uniq",
				// Creates a new instance overload on list type.
				decls.NewInstanceOverload("uniq", []*expr.Type{
					decls.NewListType(decls.Dyn),
				}, decls.Bool)),
			decls.NewVar("data", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)
	if err != nil {
		panic(err)
	}
	ast, issues := env.Compile(filter)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %v", issues.Err())
	}

	globalFunctions := cel.Functions(
		&functions.Overload{
			Operator: "uniq",
			Unary: func(lhs ref.Val) ref.Val {
				if types.ListType != lhs.Type() {
					return types.ValOrErr(lhs, "no such overload")
				}
				items, ok := lhs.Value().([]interface{})
				if !ok {
					return types.ValOrErr(lhs, "no such overload")
				}
				m := make(map[interface{}]bool)
				for _, item := range items {
					if _, exists := m[item]; exists {
						return types.Bool(false)
					}
					m[item] = true
				}
				return types.Bool(true)
			},
		},
	)
	prg, err := env.Program(ast, globalFunctions)
	if err != nil {
		log.Fatalf("program construction error: %v", err)
	}

	// Conversion from go struct -> JSON -> structpb.
	b, err := json.Marshal(myList)
	if err != nil {
		log.Fatalf("error marshalling struct: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		log.Fatalf("error unmarshalling struct: %v", err)
	}

	log.Printf("got input: %v\n", m)
	out, details, err := prg.Eval(map[string]interface{}{
		"data": m,
	})
	if err != nil {
		log.Fatalf("error evaluating program: %v", err)
	}
	log.Printf("got details: %v\n", details)
	log.Printf("got output: %v", out.Value().(bool))
}
