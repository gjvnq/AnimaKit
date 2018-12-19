package AnimaKit

import (
	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

var TheAnimation Animation

type Animation struct {
	Rect   sdl.Rect
	Length float64
	Stage  Viz
	FPS    float64
}

func (self *Animation) DrawOn(frame int, surface *sdl.Surface) error {
	return self.Stage.DrawOn(frame, surface)

}

func ffi_Animation_get_width(call otto.FunctionCall) otto.Value {
	return toValueOrPanic(TheAnimation.Rect.W)
}

func ffi_Animation_set_width(call otto.FunctionCall) otto.Value {
	tmp, err := call.Argument(1).ToInteger()
	panicOnError(err)
	TheAnimation.Rect.W = int32(tmp)
	return otto.Value{}
}

func ffi_Animation_get_height(call otto.FunctionCall) otto.Value {
	return toValueOrPanic(TheAnimation.Rect.H)
}

func ffi_Animation_set_height(call otto.FunctionCall) otto.Value {
	tmp, err := call.Argument(1).ToInteger()
	panicOnError(err)
	TheAnimation.Rect.H = int32(tmp)
	return otto.Value{}
}

func ffi_Animation_get_length(call otto.FunctionCall) otto.Value {
	return toValueOrPanic(TheAnimation.Length)
}

func ffi_Animation_set_length(call otto.FunctionCall) otto.Value {
	tmp, err := call.Argument(1).ToFloat()
	panicOnError(err)
	TheAnimation.Length = tmp
	return otto.Value{}
}

func ffi_Animation_get_fps(call otto.FunctionCall) otto.Value {
	return toValueOrPanic(TheAnimation.FPS)
}

func ffi_Animation_set_fps(call otto.FunctionCall) otto.Value {
	tmp, err := call.Argument(1).ToFloat()
	panicOnError(err)
	TheAnimation.FPS = tmp
	return otto.Value{}
}

func ffi_Animation_set_stage(call otto.FunctionCall) otto.Value {
	id, err := call.Argument(1).ToInteger()
	panicOnError(err)
	TheAnimation.Stage = mapperGet(int(id)).(Viz)
	return otto.Value{}
}
