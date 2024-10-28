package main

import (
	"fmt"

	"kcl-lang.io/kcl-go"
)

func main() {
	data := `{
		"name": "John",
		"age": 10
	}`
	code := `
	schema User:
		name: str
		age: 10

		check:
			name != "", "name is required"
			age > 10, "age must > 10"
	`
	ok, err := kcl.ValidateCode(data, code, &kcl.ValidateOptions{})
	fmt.Println(ok, err)
}
