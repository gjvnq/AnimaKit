package AnimaKit

import (
	"fmt"
	gif_lib "image/gif"
	"os"
	"sort"
	"strconv"

	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

type GIF struct {
	PositionableBase
	ScalableBase
	VisibleBase
	Frames    []*sdl.Surface
	LenFrames int
	Segs      []GifSeg
}

func NewGIFFromFile(path string) *GIF {
	path = fixPath(path)

	file, err := os.Open(path)
	panicOnError(err)

	rwops := sdl.RWFromFile(path, "r")
	format := getImgFormat(rwops)

	if format != IMG_GIF {
		panic("unimplemented format: " + strconv.Itoa(format))
	}

	gif, err := gif_lib.DecodeAll(file)
	panicOnError(err)

	ans := new(GIF)
	ans.Frames = make([]*sdl.Surface, len(gif.Image))
	for i, frame := range gif.Image {
		// Copy frame pixel by pixel to an SDL surface
		surf, err := sdl.CreateRGBSurfaceWithFormat(
			0,
			int32(frame.Rect.Dx()),
			int32(frame.Rect.Dy()),
			32,
			PIXEL_FORMAT)

		panicOnError(err)
		for x := 0; x < frame.Rect.Dx(); x++ {
			for y := 0; y < frame.Rect.Dy(); y++ {
				surf.Set(x, y, frame.At(x, y))
			}
		}
		// Store it surface
		ans.Frames[i] = surf
	}

	return ans
}

func (self GIF) Frame(frame float64) *sdl.Surface {
	for _, seg := range self.Segs {
		if seg.StartTime <= frame && frame < seg.EndTime {
			return self.Frames[seg.WhichFrame(frame)]
		}
	}
	last_seg := self.Segs[len(self.Segs)-1]
	return self.Frames[last_seg.WhichFrame(frame)]
}

func (self GIF) DrawOn(frame float64, final_surf *sdl.Surface) error {
	frame_surf := self.Frame(frame)
	// Do we even have something to draw?
	if frame_surf == nil {
		return nil
	}
	fmt.Println("Scale", self.Scale.ValAt(frame))
	fmt.Println("Visible", self.Visible.ValAt(frame))
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

type GifSeg struct {
	// Both trims are inclusive
	TrimStart int
	TrimEnd   int
	StartTime float64
	EndTime   float64
}

func (self GifSeg) WhichFrame(current_frame float64) int {
	top := self.TrimEnd - self.TrimStart
	current_frame_int := int(current_frame) - self.TrimStart
	if top == 0 {
		fmt.Println(">frame", self.TrimStart)
		return self.TrimStart + 0
	}
	ans := self.TrimStart + current_frame_int%top
	fmt.Println("frame", ans)
	return ans
}

func get_GIF(id otto.Value) *GIF {
	id_int, err := id.ToInteger()
	panicOnError(err)

	return mapperGet(int(id_int)).(*GIF)
}

func ffi_GIF_get_frames(call otto.FunctionCall) otto.Value {
	gif := get_GIF(call.Argument(0))
	return toValueOrPanic(len(gif.Frames))
}

func ffi_GIF_new(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) != 1 {
		fmt.Println(call.ArgumentList)
		panicAndPrint("Wrong number of arguments for: new GIF(filename)")
	}

	src, err := call.Argument(0).ToString()
	panicOnError(err)

	gif := NewGIFFromFile(src)

	return toValueOrPanic(mapperAdd(gif))
}

func ffi_GIF_get_keyframes(call otto.FunctionCall) otto.Value {
	return otto.Value{}
}

func ffi_GIF_set_keyframes(call otto.FunctionCall) otto.Value {
	gif := get_GIF(call.Argument(0))
	map_obj := call.Argument(1).Object()

	// Get and sort keys
	keys := make([]float64, 0)
	for _, key := range map_obj.Keys() {
		val, err := strconv.Atoi(key)
		panicOnError(err)
		keys = append(keys, float64(val))
	}
	sort.Float64s(keys)

	// Pre add first segment
	if keys[0] != 0 {
		end_time := float64(TheAnimation.Frames)
		gif.Segs = append(gif.Segs, GifSeg{
			StartTime: 0,
			EndTime:   end_time,
			TrimStart: 0,
			TrimEnd:   len(gif.Frames) - 1,
		})

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
		gif.pos_parse(keys, obj)
		gif.scale_parse(keys, obj)
		gif.visible_parse(keys, obj)
		// for _, key_int := range keys {
		// 	key := strconv.Itoa(key_int)
		// 	params := obj[key].(map[string]interface{})
		// 	i := len(gif.Segs)

		// 	start_time, err := strconv.ParseFloat(key, 64)
		// 	panicOnError(err)

		// 	new_seg := GifSeg{
		// 		StartTime: start_time,
		// 		EndTime:   end_time,
		// 	}
		// 	if params["x"] != nil {
		// 		new_seg.X = NewInterpSeg(
		// 			start_time,
		// 			end_time,
		// 			num2float64(params["x"]),
		// 			num2float64(params["x"]),
		// 			nil,
		// 		)
		// 		if i != 0 {
		// 			gif.Segs[i-1].X.EndTime = start_time
		// 			gif.Segs[i-1].X.EndVal = num2float64(params["x"])
		// 		}
		// 	}
		// 	if params["y"] != nil {
		// 		new_seg.Y = NewInterpSeg(
		// 			start_time,
		// 			end_time,
		// 			num2float64(params["y"]),
		// 			num2float64(params["y"]),
		// 			nil,
		// 		)
		// 		if i != 0 {
		// 			gif.Segs[i-1].Y.EndTime = start_time
		// 			gif.Segs[i-1].Y.EndVal = num2float64(params["y"])
		// 		}
		// 	}
		// 	if params["scale"] != nil {
		// 		new_seg.Scale = NewInterpSeg(
		// 			start_time,
		// 			end_time,
		// 			num2float64(params["scale"]),
		// 			num2float64(params["scale"]),
		// 			nil,
		// 		)
		// 		if i != 0 {
		// 			gif.Segs[i-1].Scale.EndTime = start_time
		// 			gif.Segs[i-1].Scale.EndVal = num2float64(params["scale"])
		// 		}
		// 	}

		// 	gif.Segs = append(gif.Segs, new_seg)
		// 	if i != 0 {
		// 		gif.Segs[i-1].EndTime = start_time
		// 	}
		// }
	}

	return otto.Value{}
}
