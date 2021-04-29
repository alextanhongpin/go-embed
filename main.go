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

	lv := L.Get(-1)
	if str, ok := lv.(lua.LString); ok {
		fmt.Println(string(str))
	}
	if lv.Type() != lua.LTString {
		panic("string required")
	}
}
