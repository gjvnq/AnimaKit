package AnimaKit

import (
	"fmt"

	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

type HiBitStage struct {
	Rect     sdl.Rect
	Children []Viz
	BG       ColorMixer
}

func NewHiBitStage(width, height int32) *HiBitStage {
	ans := new(HiBitStage)
	ans.Rect.W = width
	ans.Rect.H = height

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

	return toValueOrPanic(mapperAdd(NewHiBitStage(int32(width), int32(height))))
}

func ffi_HiBitStage_get_bg(call otto.FunctionCall) otto.Value {
	stage := get_HiBitStage(call.Argument(0))

	return toValueOrPanic(stage.BG)
}

func ffi_HiBitStage_set_bg(call otto.FunctionCall) otto.Value {
	stage := get_HiBitStage(call.Argument(0))

	stage.BG.FromValue(call.Argument(1))
	return otto.Value{}
}

func (self *HiBitStage) DrawOn(frame int, final_surf *sdl.Surface) error {
	// Create surface of output size
	virtual_surf, err := sdl.CreateRGBSurfaceWithFormat(0, self.Rect.W, self.Rect.H, 32, PIXEL_FORMAT)
	panicOnError(err)

	// Draw background
	virtual_surf.FillRect(&self.Rect, color2uint32(self.BG.ValAt(float64(frame))))

	// Copy to output
	r := RectFitAndCenterInSurf(self.Rect, final_surf)
	virtual_surf.BlitScaled(nil, final_surf, r)

	return nil
}
