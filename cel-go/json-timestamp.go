package main

import (
	"log"
	"time"

	//structpb "google.golang.org/protobuf/types/known/structpb"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

// An example of using native JSON-friendly Go struct as a type and input in a Cel program.
func main() {

	// NOTE: myStruct.number == 10 will fail with error: no such overload.
	//https://github.com/google/cel-go/issues/203
	filter := `date > now`

	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("date", decls.Timestamp),
			decls.NewVar("now", decls.Timestamp),
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

	m := map[string]interface{}{
		"date": time.Now().Add(5 * time.Second),
		"now":  time.Now(),
	}
	log.Printf("got input: %v\n", m)
	out, details, err := prg.Eval(m)
	if err != nil {
		log.Fatalf("error evaluating program: %v", err)
	}
	log.Printf("got details: %v\n", details)
	log.Printf("got output: %v", out.Value().(bool))
}
