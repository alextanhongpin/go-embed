package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type User struct {
	Name string `json:"name"`
}

type Env struct {
	Users []User `json:"users"`
}

func main() {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewFunction("lowercase", decls.NewOverload("lowercase_string", []*expr.Type{decls.String}, decls.String)),
			decls.NewVar("users", decls.NewListType(decls.Dyn)),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	globalFunctions := cel.Functions(
		&functions.Overload{
			Operator: "lowercase_string",
			Unary: func(lhs ref.Val) ref.Val {
				return types.String(strings.ToLower(lhs.Value().(string)))
			},
		},
	)

	ast, issues := env.Compile(`users.all(u, lowercase(u.name).startsWith("a"))`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	prg, err := env.Program(ast, globalFunctions)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	m, err := structToMap(Env{
		Users: []User{
			User{Name: "Andrew"},
			User{Name: "Alpha"},
			User{Name: "aloha"},
			//User{Name: "Carl"},
		},
	})
	if err != nil {
		panic(err)
	}
	out, _, err := prg.Eval(m)
	if err != nil {
		log.Fatalf("program err: %s", err)
	}
	log.Printf("got output: %v", out)
}

func structToMap(v interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
