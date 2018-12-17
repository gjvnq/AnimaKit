package main

import (
	"fmt"

	"github.com/gjvnq/AnimaKit"
)

func main() {
	vm, err := AnimaKit.LoadScriptFromFile("examples/typewriter.js")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	value, err := vm.Call("prepareStage", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(value)
}
