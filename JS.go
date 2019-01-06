package AnimaKit

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/robertkrimen/otto"
)

var VM = otto.New()
var ScriptFolder = ""

func fixPath(path string) string {
	return filepath.Join(ScriptFolder, path)
}

func LoadScriptFromFile(path string) (*otto.Otto, error) {
	path, err := filepath.Abs(path)
	panicOnError(err)
	ScriptFolder = filepath.Dir(path)

	// Load file
	TheLog.InfoF("Loading animation script from: %s", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Add functions
	VM.Set("println", JS_Println)
	VM.Set("ffi_HiBitStage_new", ffi_HiBitStage_new)
	VM.Set("ffi_HiBitStage_get_bg", ffi_HiBitStage_get_bg)
	VM.Set("ffi_HiBitStage_set_bg", ffi_HiBitStage_set_bg)
	VM.Set("ffi_HiBitStage_place", ffi_HiBitStage_place)

	VM.Set("ffi_Animation_get_width", ffi_Animation_get_width)
	VM.Set("ffi_Animation_set_width", ffi_Animation_set_width)
	VM.Set("ffi_Animation_get_height", ffi_Animation_get_height)
	VM.Set("ffi_Animation_set_height", ffi_Animation_set_height)
	VM.Set("ffi_Animation_get_length", ffi_Animation_get_length)
	VM.Set("ffi_Animation_set_length", ffi_Animation_set_length)
	VM.Set("ffi_Animation_get_fps", ffi_Animation_get_fps)
	VM.Set("ffi_Animation_set_fps", ffi_Animation_set_fps)
	VM.Set("ffi_Animation_set_stage", ffi_Animation_set_stage)

	VM.Set("ffi_GIF_new", ffi_GIF_new)
	VM.Set("ffi_GIF_get_frames", ffi_GIF_get_frames)
	VM.Set("ffi_GIF_get_keyframes", ffi_GIF_get_keyframes)
	VM.Set("ffi_GIF_set_keyframes", ffi_GIF_set_keyframes)

	VM.Set("ffi_Image_new", ffi_Image_new)
	VM.Set("ffi_Image_get_keyframes", ffi_Image_get_keyframes)
	VM.Set("ffi_Image_set_keyframes", ffi_Image_set_keyframes)

	// Load wrapper
	scripts, _ := AssetDir("res")
	for _, script := range scripts {
		// Only load .js files
		if !strings.HasSuffix(script, ".js") {
			continue
		}
		file, err := Asset("res/" + script)
		panicOnError(err)
		_, err = VM.Run(file)
		panicOnError(err)
	}

	// Execute it
	TheLog.InfoF("Running animation script from: %s", path)
	_, err = VM.Run(file)
	panicOnError(err)
	TheLog.InfoF("[FINISHED] Ran animation script from: %s", path)

	return VM, nil
}

func jsval2string(val otto.Value) string {
	ans := ""

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
