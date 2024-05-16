package main

import (
	"github/mdtosif/consumet/gogoanime"
	"syscall/js"
)

func add(this js.Value, i []js.Value) interface{} {
	arg1 := i[0].String()
	arg2 := i[0].Int()
	data, err := gogoanime.Search(arg1, arg2)
	if err != nil {
		return err.Error()
	}
	return data
}

func main() {
	// Register the add function to be callable from JavaScript
	js.Global().Set("add", js.FuncOf(add))

	// Keep the program running indefinitely
	select {}
}
