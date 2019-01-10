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
	src       string
}

func NewGIFFromFile(path string) *GIF {
	path = fixPath(path)

	file, err := os.Open(path)
	defer file.Close()
	panicOnError(err)

	// Ensure auto reload
	AddFileToWatch(path)

	rwops := sdl.RWFromFile(path, "r")
	format := getImgFormat(rwops)

	if format != IMG_GIF {
		panic("image is not a GIF: " + path)
	}

	gif, err := gif_lib.DecodeAll(file)
	panicOnError(err)

	ans := new(GIF)
	ans.src = path
	ans.Frames = make([]*sdl.Surface, len(gif.Image))
	for i, frame := range gif.Image {
		// Prepare surface
		surf, err := sdl.CreateRGBSurfaceWithFormat(
			0,
			int32(frame.Rect.Dx()),
			int32(frame.Rect.Dy()),
			32,
			PIXEL_FORMAT)
		panicOnError(err)

		// Copy frame pixel by pixel to an SDL surface
		pix := surf.Pixels()
		for x := 0; x < frame.Rect.Dx(); x++ {
			for y := 0; y < frame.Rect.Dy(); y++ {
				i := int32(y)*surf.Pitch + int32(x)*int32(surf.Format.BytesPerPixel)
				r, g, b, a := frame.At(x, y).RGBA()
				pix[i+3] = limit_byte(int(a / 0xff))
				pix[i+2] = limit_byte(int(r / 0xff))
				pix[i+1] = limit_byte(int(g / 0xff))
				pix[i+0] = limit_byte(int(b / 0xff))
			}
		}

		// Store it surface
		ans.Frames[i] = surf
	}

	ans.pos_init()
	ans.scale_init()
	ans.visible_init()

	return ans
}

func (self GIF) Frame(frame float64) (*sdl.Surface, int) {
	for _, seg := range self.Segs {
		if seg.StartTime <= frame && frame < seg.EndTime {
			ans := seg.WhichFrame(frame)
			return self.Frames[ans], ans
		}
	}
	if len(self.Segs) > 0 {
		last_seg := self.Segs[len(self.Segs)-1]
		ans := last_seg.WhichFrame(frame)
		return self.Frames[ans], ans
	}
	return self.Frames[0], 0
}

func (self GIF) DrawOn(frame float64, final_surf *sdl.Surface) error {
	frame_surf, _ := self.Frame(frame)
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

type GifSeg struct {
	// Both trims are inclusive
	TrimStart int
	TrimStop  int
	StartTime float64
	EndTime   float64
	Delay     float64 //in frames
}

func (self GifSeg) WhichFrame(current_frame float64) int {
	top := self.TrimStop - self.TrimStart
	current_frame_int := int(current_frame) - self.TrimStart
	if top == 0 {
		return self.TrimStart + 0
	}
	ans := self.TrimStart + int(float64(current_frame_int)/self.Delay)%top

	// Safety
	if ans < self.TrimStart {
		ans = self.TrimStart
	}
	if ans > self.TrimStop {
		ans = self.TrimStop
	}

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

	TheLog.DebugF("Setting keyframes for GIF(%s)", gif.src)
	defer TheLog.DebugF("[FINISHED] Setting keyframes for GIF(%s)", gif.src)

	// Get and sort keys
	sorted_keys := make([]float64, 0)
	for _, key := range map_obj.Keys() {
		val, err := strconv.Atoi(key)
		panicOnError(err)
		sorted_keys = append(sorted_keys, float64(val))
	}
	sort.Float64s(sorted_keys)

	// Pre add first segment
	if sorted_keys[0] != 0 {
		end_time := float64(TheAnimation.Frames)
		gif.Segs = append(gif.Segs, GifSeg{
			StartTime: 0,
			EndTime:   end_time,
			TrimStart: 0,
			TrimStop:  len(gif.Frames) - 1,
		})
	}

	// Get values for each key frame
	tmp, err := call.Argument(1).Export()
	panicOnError(err)
	key_frame_spec := make(map[float64]map[string]interface{})
	for key, val := range tmp.(map[string]interface{}) {
		key_int, err := strconv.Atoi(key)
		panicOnError(err)
		tmp2, ok := val.(map[string]interface{})
		if ok {
			key_frame_spec[float64(key_int)] = tmp2
		}
	}
	// Go in order
	gif.pos_parse(sorted_keys, key_frame_spec)
	gif.scale_parse(sorted_keys, key_frame_spec)
	gif.visible_parse(sorted_keys, key_frame_spec)

	// Ensure default value
	has_defualt := false
	if sorted_keys[0] == 0 {
		_, has_defualt = key_frame_spec[0]["gif"]
	}
	gif.Segs = make([]GifSeg, 0)
	if !has_defualt {
		gif.Segs = make([]GifSeg, 1)
		gif.Segs = append(gif.Segs, GifSeg{
			TrimStart: 0,
			TrimStop:  len(gif.Frames),
			StartTime: 0,
			EndTime:   PosInf,
			Delay:     0.1 * TheAnimation.FPS,
		})
	}
	for _, key := range sorted_keys {
		params := key_frame_spec[key]

		// Ignore key frames that aren't about us
		if _, ok := params["gif"]; !ok {
			continue
		}

		// Load info
		params = params["gif"].(map[string]interface{})

		trim_start := 0
		if val, ok := params["trim_start"]; ok {
			trim_start = num2int(val)
		}
		trim_stop := len(gif.Frames)
		if val, ok := params["trim_stop"]; ok {
			trim_stop = num2int(val)
		}
		delay := 0.1 * TheAnimation.FPS
		if val, ok := params["delay"]; ok {
			delay = num2float64(val)
		}

		// Safety
		if delay <= 0 {
			delay = PosInf
		}

		// Add key frame
		if len(gif.Segs) > 1 {
			gif.Segs[len(gif.Segs)-1].EndTime = key
		}
		gif.Segs = append(gif.Segs, GifSeg{
			TrimStart: trim_start,
			TrimStop:  trim_stop,
			StartTime: key,
			EndTime:   PosInf,
			Delay:     delay,
		})
		TheLog.DebugF("at frame %0.f: %#+v", key, gif.Segs[len(gif.Segs)-1])
	}
	TheLog.DebugF("gif.Segs: %#+v", gif.Segs)

	return otto.Value{}
}
