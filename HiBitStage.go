package AnimaKit

import (
	"image/color"

	"github.com/robertkrimen/otto"
)

var MapperHiBitStage = make([]*HiBitStage, 0)

type HiBitStage struct {
	Rect     Rect
	FPS      float64
	Children []Viz
	BG       color.NRGBA
}

func NewHiBitStage(width, height int, fps float64) *HiBitStage {
	ans := new(HiBitStage)
	ans.Rect.Width = width
	ans.Rect.Height = height
	ans.FPS = fps
	ans.BG = color.NRGBA{0, 0, 0, 255}

	return ans
}

func get_HiBitStage(id otto.Value) *HiBitStage {
	id_int, err := id.ToInteger()
	panicOnError(err)

	return MapperHiBitStage[int(id_int)]
}

func ffi_HiBitStage_new(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 3 {
		panicAndPrint("Wrong number of arguments for: new HiBitStage(width, height, fps)")
	}

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

	return toValueOrPanic(NRGBA2hex(stage.BG))
}

func ffi_HiBitStage_set_bg(call otto.FunctionCall) otto.Value {
	stage := get_HiBitStage(call.Argument(0))

	var err error
	hex, err := call.Argument(1).ToString()
	panicOnError(err)
	stage.BG = hex2NRGBA(hex)
	return otto.Value{}
}
