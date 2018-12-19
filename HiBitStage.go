package AnimaKit

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

type HiBitStage struct {
	Rect     Rect
	Children []Viz
	BG       color.NRGBA
}

func NewHiBitStage(width, height int) *HiBitStage {
	ans := new(HiBitStage)
	ans.Rect.Width = width
	ans.Rect.Height = height
	ans.BG = color.NRGBA{0, 0, 0, 255}

	return ans
}

func get_HiBitStage(id otto.Value) *HiBitStage {
	id_int, err := id.ToInteger()
	panicOnError(err)

	return mapperGet(int(id_int)).(*HiBitStage)
}

func ffi_HiBitStage_new(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 2 {
		fmt.Println(call.ArgumentList)
		panicAndPrint("Wrong number of arguments for: new HiBitStage(width, height)")
	}

	width, err := call.Argument(0).ToInteger()
	panicOnError(err)
	height, err := call.Argument(1).ToInteger()
	panicOnError(err)

	return toValueOrPanic(mapperAdd(NewHiBitStage(int(width), int(height))))
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

func (self *HiBitStage) DrawOn(frame int, surface sdl.Surface) error {
	return errors.New("not implemented")
}
