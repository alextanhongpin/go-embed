package main

import (
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func main() {
	cel.EvalOptions(cel.OptTrackState)

	// Expose name and group variables to CEL using the
	// cel.Declarations environment option.
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("name", decls.String),
			decls.NewVar("group", decls.String),
		),
	)

	// Parsing checks if the expression is syntactically
	// valid.
	ast, issues := env.Compile(`name.startsWith("/groups/" + group)`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	// The cel.Program generated at the end of the parse
	// and check is stateless, thread-safe and cachable.
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	// The evaluation is thread-safe and side-effect free.
	// Fields that are present in the input, but not referenced in the expression are ignored.
	out, details, err := prg.Eval(map[string]interface{}{
		"name":  "/groups/acme.co/documents/secret-stuff",
		"group": "acme.co",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(details)
	log.Println(out)
}
