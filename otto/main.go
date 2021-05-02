package main

import (
	"fmt"
	"log"

	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()
	_, err := vm.Run(`
		// const, let, arrow functions not supported...
		var nums = [1,2,3]
		var a = nums.reduce(function (a, b) {
			return a + b
		}, 0)
		abc = 2 + a;
		console.log("The value of abc is " + abc) // 4
	`)
	if err != nil {
		panic(err)
	}
	// Get a value out of the VM.
	if value, err := vm.Get("abc"); err == nil {
		if v, err := value.ToInteger(); err == nil {
			fmt.Println("got", v)
		}
	}

	// Set a number.
	vm.Set("def", 11)
	vm.Run(`
		console.log("The value of def is " + def); // 11
	`)

	// Set a string.
	vm.Set("xyz", "Nothing happens.")
	vm.Run(`
		console.log(xyz.length)
	`)

	// Get the value of an expression.
	value, _ := vm.Run("xyz.length")
	{
		v, _ := value.ToInteger()
		fmt.Println("expression value is", v)
	}

	// Error.
	value, err = vm.Run("helloworld.length")
	if err != nil {
		log.Println(err)
		fmt.Println(value.IsUndefined()) // true
	}

	// Set a go function.

	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s", call.Argument(0).String())
		return otto.Value{}
	})

	// Set a go function that results something usefule.
	vm.Set("twoPlus", func(call otto.FunctionCall) otto.Value {
		right, _ := call.Argument(0).ToInteger()
		result, _ := vm.ToValue(2 + right)
		return result
	})

	// Use the functions in JavaScript.
	value, err = vm.Run(`
		sayHello("World"); // Hello, World
		sayHello(); // Hello, undefined

		result = twoPlus(2.0); // 4
	`)
	if err != nil {
		log.Fatal(err)
	}
	{
		v, _ := value.ToInteger()
		log.Println("received", v)
	}
}
