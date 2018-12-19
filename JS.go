package AnimaKit

import (
	"fmt"
	"os"

	"github.com/robertkrimen/otto"
)

var VM = otto.New()

func LoadScriptFromFile(path string) (*otto.Otto, error) {
	// Load file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Add functions
	VM.Set("println", JS_Println)
	VM.Set("ffi_HiBitStage_new", ffi_HiBitStage_new)
	VM.Set("ffi_HiBitStage_get_bg", ffi_HiBitStage_get_bg)
	VM.Set("ffi_HiBitStage_set_bg", ffi_HiBitStage_set_bg)

	VM.Set("ffi_Animation_get_width", ffi_Animation_get_width)
	VM.Set("ffi_Animation_set_width", ffi_Animation_set_width)
	VM.Set("ffi_Animation_get_height", ffi_Animation_get_height)
	VM.Set("ffi_Animation_set_height", ffi_Animation_set_height)
	VM.Set("ffi_Animation_get_length", ffi_Animation_get_length)
	VM.Set("ffi_Animation_set_length", ffi_Animation_set_length)
	VM.Set("ffi_Animation_get_fps", ffi_Animation_get_fps)
	VM.Set("ffi_Animation_set_fps", ffi_Animation_set_fps)
	VM.Set("ffi_Animation_set_stage", ffi_Animation_set_stage)

	// Load wrapper
	scripts := []string{"res/HiBitStage.js", "res/Animation.js"}
	for _, script := range scripts {
		file, err := Asset(script)
		panicOnError(err)
		_, err = VM.Run(file)
		panicOnError(err)
	}

	// Execute it
	_, err = VM.Run(file)
	panicOnError(err)

	return VM, nil
}

func jsval2string(val otto.Value) string {
	ans := ""
	// exp, err := val.Export()
	// panicOnError(err)

	// switch exp_val := exp.(type) {
	// case string:
	// 	return "\"" + exp_val + "\""
	// case map[string]interface{}:
	// case []interface{}:
	// 	ans = "["
	// 	for _, item := range exp_val {
	// 		if i != 0 {
	// 			ans += ", "
	// 		}
	// 		ans += key
	// 	}
	// 	ans += "]"
	// 	return ans
	// default:
	// 	return val.String()
	// }

	if val.Class() == "Object" {
		ans = "Object{"
		for i, key := range val.Object().Keys() {
			if i != 0 {
				ans += ", "
			}
			ans += key
		}
		ans += "}"
		return ans
	} else if val.Class() == "Array" {
		ans = "["
		for i, key := range val.Object().Keys() {
			if i != 0 {
				ans += ", "
			}
			subval, err := val.Object().Get(key)
			panicOnError(err)
			ans += jsval2string(subval)
		}
		ans += "]"
		return ans
	} else {
		return val.String()
	}

}

func JS_Println(call otto.FunctionCall) otto.Value {
	fmt.Printf("JS: ")

	for i, arg := range call.ArgumentList {
		if i != 0 {
			fmt.Printf(" ")
		}
		fmt.Printf(jsval2string(arg))
	}
	fmt.Printf("\n")

	return otto.Value{}
}
