package main

import (
	"fmt"
	"time"

	"github.com/antonmedv/expr"
)

type Env struct {
	Tweets []Tweet
}

func (Env) Format(t time.Time) string {
	return t.Format(time.RFC822)
}

type Tweet struct {
	Text string
	Date time.Time
}

func main() {
	code := `map(filter(Tweets, {len(.Text) > 0}), {.Text + Format(.Date)})`

	// We can use an empty instance of the struct as a n environment.
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		panic(err)
	}
	env := Env{
		Tweets: []Tweet{
			{"Oh My God!", time.Now()},
			{"How you doing?", time.Now()},
			{"Could I be wearing any more clothes?", time.Now()},
			{"", time.Now()},
		},
	}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
	result, ok := output.([]interface{})
	if !ok {
		panic("invalid type casting")
	}
	for _, r := range result {
		fmt.Println(r.(string))
	}
}
