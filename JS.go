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

	// Load wrapper
	scripts := []string{"res/HiBitStage.js"}
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

func JS_Println(call otto.FunctionCall) otto.Value {
	fmt.Println("JS: ", call.ArgumentList)
	return otto.Value{}
}
