package AnimaKit

import (
	"fmt"
	gif_lib "image/gif"
	"os"
	"strconv"

	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

type GIF struct {
	Frames    []*sdl.Surface
	LenFrames int
	Segs      []CyclableSegments
}

func NewGIFFromFile(path string) GIF {
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

	ans := GIF{}
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
	for _, seg := range self.Segs {
		if seg.StartFrame <= frame && frame < seg.StopFrame {
			return self.Frames[seg.WhichFrame(frame)]
		}
	}
	return nil
}

type CyclableSegments struct {
	StartFrame int
	StopFrame  int
	// Both trims are inclusive
	TrimStart int
	TrimEnd   int
	Speed     float64
	Cycle     bool
}

func (self CyclableSegments) WhichFrame(current_frame int) int {
	current_frame -= self.StartFrame

	frame_walk := float64(current_frame) * self.Speed
	out_frame := current_frame + self.TrimStart + int(frame_walk)
	self_len := self.TrimEnd - self.TrimStart
	if !self.Cycle {
		return out_frame
	}
	for out_frame > self.TrimEnd {
		out_frame -= self_len
	}
	for out_frame < self.TrimStart {
		out_frame += self_len
	}
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

	return toValueOrPanic(mapperAdd(NewGIFFromFile(src)))
}

func ffi_GIF_get_keyframes(call otto.FunctionCall) otto.Value {
	return otto.Value{}
}

func ffi_GIF_set_keyframes(call otto.FunctionCall) otto.Value {
	return otto.Value{}
}
