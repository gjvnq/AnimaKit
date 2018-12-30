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

func (self GIF) Frame(frame int) *sdl.Surface {
	frame_f64 := float64(frame)

	for _, seg := range self.Segs {
		if seg.StartTime <= frame_f64 && frame_f64 < seg.EndTime {
			fmt.Println(seg.WhichFrame(frame))
			return self.Frames[seg.WhichFrame(frame)]
		}
	}
	return nil
}

func (self GIF) DrawOn(frame int, final_surf *sdl.Surface) error {
	frame_surf := self.Frame(frame)
	if frame_surf == nil {
		return nil
	}

	// r := RectFitAndCenterInSurf(frame_surf.Rect, final_surf)
	return frame_surf.BlitScaled(nil, final_surf, nil)
}

type GifSeg struct {
	// Both trims are inclusive
	TrimStart int
	TrimEnd   int
	Visible   bool
	StartTime float64
	EndTime   float64
	Speed     InterpSeg
	Cycle     bool
	X         InterpSeg
	Y         InterpSeg
	Scale     InterpSeg
}

func (self GifSeg) WhichFrame(current_frame int) int {
	fmt.Println("->", current_frame)

	current_frame -= int(self.StartTime)

	frame_walk := float64(current_frame) * self.Speed.ValAtᐸintᐳ(current_frame)
	out_frame := current_frame + self.TrimStart + int(frame_walk)
	self_len := self.TrimEnd - self.TrimStart
	if !self.Cycle {
		fmt.Println(">", out_frame)
		return out_frame
	}
	for out_frame > self.TrimEnd {
		out_frame -= self_len
	}
	for out_frame < self.TrimStart {
		out_frame += self_len
	}

	fmt.Println(">", out_frame)
	return out_frame
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
	keys := make([]int, 0)
	for _, key := range map_obj.Keys() {
		val, err := strconv.Atoi(key)
		panicOnError(err)
		keys = append(keys, val)
	}
	sort.Ints(keys)

	// Pre add first segment
	if keys[0] != 0 {
		end_time := float64(TheAnimation.Frames)
		gif.Segs = append(gif.Segs, GifSeg{
			StartTime: 0,
			EndTime:   end_time,
			Visible:   false,
		})

		// Get values for each key frame
		tmp, err := call.Argument(1).Export()
		panicOnError(err)
		obj := tmp.(map[string]interface{})
		// Go in order
		for _, key_int := range keys {
			key := strconv.Itoa(key_int)
			params := obj[key].(map[string]interface{})
			i := len(gif.Segs)

			start_time, err := strconv.ParseFloat(key, 64)
			panicOnError(err)

			new_seg := GifSeg{
				StartTime: start_time,
				EndTime:   end_time,
			}
			if params["x"] != nil {
				new_seg.X = NewInterpSeg(
					start_time,
					end_time,
					num2float64(params["x"]),
					num2float64(params["x"]),
					nil,
				)
				if i != 0 {
					gif.Segs[i-1].X.EndTime = start_time
					gif.Segs[i-1].X.EndVal = num2float64(params["x"])
				}
			}
			if params["y"] != nil {
				new_seg.Y = NewInterpSeg(
					start_time,
					end_time,
					num2float64(params["y"]),
					num2float64(params["y"]),
					nil,
				)
				if i != 0 {
					gif.Segs[i-1].Y.EndTime = start_time
					gif.Segs[i-1].Y.EndVal = num2float64(params["y"])
				}
			}
			if params["scale"] != nil {
				new_seg.Scale = NewInterpSeg(
					start_time,
					end_time,
					num2float64(params["scale"]),
					num2float64(params["scale"]),
					nil,
				)
				if i != 0 {
					gif.Segs[i-1].Scale.EndTime = start_time
					gif.Segs[i-1].Scale.EndVal = num2float64(params["scale"])
				}
			}

			gif.Segs = append(gif.Segs, new_seg)
			if i != 0 {
				gif.Segs[i-1].EndTime = start_time
			}
		}
	}

	return otto.Value{}
}
