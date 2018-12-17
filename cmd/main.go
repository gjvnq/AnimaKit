package main

import (
	"fmt"

	"github.com/gjvnq/AnimaKit"
	"github.com/integrii/flaggy"
)

var ArgsRes string = ""
var ArgsInput string = ""
var ArgsOutput string = ""

func init() {
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("AnimaKit")
	flaggy.SetDescription("A simple animation renderer")
	flaggy.DefaultParser.ShowHelpOnUnexpected = false
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/gjvnq/AnimaKit"
	flaggy.String(&ArgsRes, "", "res", "Output resolution. Ex: --res 800x600")
	flaggy.AddPositionalValue(&ArgsInput, "input", 1, true, ".js file with the animation")
	flaggy.AddPositionalValue(&ArgsOutput, "output", 2, false, "video output")
	flaggy.Parse()
}

func main() {
	_, err := AnimaKit.LoadScriptFromFile(ArgsInput)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// value, err := vm.Call("prepareStage", nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(value)
}
