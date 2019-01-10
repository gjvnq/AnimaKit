package AnimaKit

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"

	sdlImg "github.com/veandco/go-sdl2/img"
)

const (
	IMG_UNKNOWN = 0
	IMG_BMP     = 1
	IMG_CUR     = 2
	IMG_GIF     = 3
	IMG_ICO     = 4
	IMG_JPG     = 5
	IMG_LBM     = 6
	IMG_PCX     = 7
	IMG_PNG     = 8
	IMG_PNM     = 9
	IMG_TIF     = 10
	IMG_WEBP    = 11
	IMG_XCF     = 12
	IMG_XV      = 13
)

func getImgFormat(src *sdl.RWops) int {
	switch {
	case sdlImg.IsBMP(src):
		return IMG_BMP
	case sdlImg.IsCUR(src):
		return IMG_CUR
	case sdlImg.IsGIF(src):
		return IMG_GIF
	case sdlImg.IsICO(src):
		return IMG_ICO
	case sdlImg.IsJPG(src):
		return IMG_JPG
	case sdlImg.IsLBM(src):
		return IMG_LBM
	case sdlImg.IsPCX(src):
		return IMG_PCX
	case sdlImg.IsPNG(src):
		return IMG_PNG
	case sdlImg.IsPNM(src):
		return IMG_PNM
	case sdlImg.IsTIF(src):
		return IMG_TIF
	case sdlImg.IsWEBP(src):
		return IMG_WEBP
	case sdlImg.IsXCF(src):
		return IMG_XCF
	case sdlImg.IsXCF(src):
		return IMG_XCF
	case sdlImg.IsXV(src):
		return IMG_XV
	default:
		return IMG_UNKNOWN
	}
}

type Image struct {
	PositionableBase
	ScalableBase
	VisibleBase
	Surface *sdl.Surface
	src     string
}

func NewImageFromFile(path string) *Image {
	path = fixPath(path)

	file, err := os.Open(path)
	defer file.Close()
	panicOnError(err)

	// Ensure auto reload
	AddFileToWatch(path)

	rwops := sdl.RWFromFile(path, "r")
	ans := new(Image)
	ans.src = path
	err = nil
	switch getImgFormat(rwops) {
	case IMG_GIF:
		ans.Surface, err = sdlImg.LoadGIFRW(rwops)
	case IMG_BMP:
		ans.Surface, err = sdlImg.LoadBMPRW(rwops)
	case IMG_PNG:
		ans.Surface, err = sdlImg.LoadPNGRW(rwops)
	case IMG_JPG:
		ans.Surface, err = sdlImg.LoadJPGRW(rwops)
	default:
		panic("Unsuported image format: " + path)
	}
	panicOnError(err)

	ans.pos_init()
	ans.scale_init()
	ans.visible_init()

	return ans
}

func (self Image) DrawOn(frame float64, final_surf *sdl.Surface) error {
	frame_surf := self.Surface
	// Do we even have something to draw?
	if frame_surf == nil {
		return nil
	}
	// Do we really need to draw?
	if !self.Visible.ValAt(frame) {
		return nil
	}

	out_rect := RectScale(Surface2Rect(frame_surf), self.Scale.ValAt(frame))
	out_rect.X = int32(self.X.ValAt(frame))
	out_rect.Y = int32(self.Y.ValAt(frame))
	out_rect = UserPos2SDLPos(out_rect, final_surf)
	return frame_surf.BlitScaled(nil, final_surf, out_rect)
}

func get_Image(id otto.Value) *Image {
	id_int, err := id.ToInteger()
	panicOnError(err)

	return mapperGet(int(id_int)).(*Image)
}

func ffi_Image_new(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		fmt.Println(call.ArgumentList)
		panicAndPrint("Wrong number of arguments for: new Image(filename)")
	}

	src, err := call.Argument(0).ToString()
	panicOnError(err)

	img := NewImageFromFile(src)

	return toValueOrPanic(mapperAdd(img))
}

func ffi_Image_get_keyframes(call otto.FunctionCall) otto.Value {
	return otto.Value{}
}

func ffi_Image_set_keyframes(call otto.FunctionCall) otto.Value {
	img := get_Image(call.Argument(0))
	map_obj := call.Argument(1).Object()

	TheLog.DebugF("Setting keyframes for Image(%s)", img.src)
	defer TheLog.DebugF("[FINISHED] Setting for Image(%s)", img.src)

	// Get and sort keys
	keys := make([]float64, 0)
	for _, key := range map_obj.Keys() {
		val, err := strconv.Atoi(key)
		panicOnError(err)
		keys = append(keys, float64(val))
	}
	sort.Float64s(keys)

	// Get values for each key frame
	tmp, err := call.Argument(1).Export()
	panicOnError(err)
	obj := make(map[float64]map[string]interface{})
	for key, val := range tmp.(map[string]interface{}) {
		key_int, err := strconv.Atoi(key)
		panicOnError(err)
		tmp2, ok := val.(map[string]interface{})
		if ok {
			obj[float64(key_int)] = tmp2
		}
	}
	// Go in order
	img.pos_parse(keys, obj)
	img.scale_parse(keys, obj)
	img.visible_parse(keys, obj)

	return otto.Value{}
}
