package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}

	if err := L.DoString(`return "this is amazing"`); err != nil {
		panic(err)
	}

	lv := L.Get(-1)
	if str, ok := lv.(lua.LString); ok {
		fmt.Println("got: " + string(str))
	}
}
