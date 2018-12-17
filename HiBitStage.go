package AnimaKit

import (
	"github.com/robertkrimen/otto"
)

var MapperHiBitStage = make([]*HiBitStage, 0)

type HiBitStage struct {
	Rect     Rect
	FPS      float64
	Children []Viz
	BG       string
}

func NewHiBitStage(width, height int, fps float64) *HiBitStage {
	ans := new(HiBitStage)
	ans.Rect.Width = width
	ans.Rect.Height = height
	ans.FPS = fps
	ans.BG = "#000000ff"

	return ans
}

func get_HiBitStage(id otto.Value) *HiBitStage {
	id_int, err := id.ToInteger()
	panicOnError(err)

	return MapperHiBitStage[int(id_int)]
}

func ffi_HiBitStage_new(call otto.FunctionCall) otto.Value {
	width, err := call.Argument(0).ToInteger()
	panicOnError(err)
	height, err := call.Argument(1).ToInteger()
	panicOnError(err)
	fps, err := call.Argument(2).ToFloat()
	panicOnError(err)

	id := len(MapperHiBitStage)
	MapperHiBitStage = append(MapperHiBitStage, NewHiBitStage(int(width), int(height), fps))

	return toValueOrPanic(id)
}

func ffi_HiBitStage_get_bg(call otto.FunctionCall) otto.Value {
	stage := get_HiBitStage(call.Argument(0))

	return toValueOrPanic(stage.BG)
}

func ffi_HiBitStage_set_bg(call otto.FunctionCall) otto.Value {
	stage := get_HiBitStage(call.Argument(0))

	var err error
	stage.BG, err = call.Argument(1).ToString()
	panicOnError(err)
	return otto.Value{}
}
