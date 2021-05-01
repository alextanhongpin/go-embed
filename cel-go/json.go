package main

import (
	"encoding/json"
	"log"

	//structpb "google.golang.org/protobuf/types/known/structpb"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

type Payload struct {
	Strs []string          `json:"strs"`
	Data map[string]string `json:"data"`
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
	filter := `int(myStruct.num) == 10 && myStruct.str == "hello" && "world" in myStruct.payload.data && "banana" in myStruct.payload.strs`
	myStruct := MyStruct{
		Num: 10,
		Str: "hello",
		Payload: Payload{
			Data: map[string]string{"world": "foobar"},
			Strs: []string{"banana"},
		},
	}

	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("myStruct", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)
	if err != nil {
		panic(err)
	}
	ast, issues := env.Compile(filter)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %v", issues.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %v", err)
	}

	// Conversion from go struct -> JSON -> structpb.
	b, err := json.Marshal(myStruct)
	if err != nil {
		log.Fatalf("error marshalling struct: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		log.Fatalf("error unmarshalling struct: %v", err)
	}

	log.Printf("got input: %v\n", m)
	out, details, err := prg.Eval(map[string]interface{}{
		"myStruct": m,
	})
	if err != nil {
		log.Fatalf("error evaluating program: %v", err)
	}
	log.Printf("got details: %v\n", details)
	log.Printf("got output: %v", out.Value().(bool))
}
