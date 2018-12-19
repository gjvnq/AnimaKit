package AnimaKit

import "github.com/robertkrimen/otto"

var TheAnimation Animation

type Animation struct {
	Rect   Rect
	Length float64
	Stage  Viz
	FPS    float64
}

func ffi_Animation_get_width(call otto.FunctionCall) otto.Value {
	return toValueOrPanic(TheAnimation.Rect.Width)
}

func ffi_Animation_set_width(call otto.FunctionCall) otto.Value {
	tmp, err := call.Argument(1).ToInteger()
	panicOnError(err)
	TheAnimation.Rect.Width = int(tmp)
	return otto.Value{}
}

func ffi_Animation_get_height(call otto.FunctionCall) otto.Value {
	return toValueOrPanic(TheAnimation.Rect.Height)
}

func ffi_Animation_set_height(call otto.FunctionCall) otto.Value {
	tmp, err := call.Argument(1).ToInteger()
	panicOnError(err)
	TheAnimation.Rect.Height = int(tmp)
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
	tmp, err := call.Argument(1).ToInteger()
	panicOnError(err)
	TheAnimation.Stage = mapperGet(int(tmp)).(Viz)
	return otto.Value{}
}
