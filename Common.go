package AnimaKit

import (
	"fmt"
	"image/color"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"time"

	logger "github.com/gjvnq/go-logger"
	"github.com/robertkrimen/otto"
	sdlImg "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const VERSION = "0.1.2-dev"

// const PIXEL_FORMAT = sdl.PIXELFORMAT_ARGB8888
const PIXEL_FORMAT = sdl.PIXELFORMAT_ARGB8888

var TheLog *logger.Logger

var PosInf = math.Inf(1)

func init() {
	var err error
	// Prepare logger
	TheLog, err = logger.New("AnimaKit", 1)
	if err != nil {
		panic(err)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicAndPrint(any interface{}) {
	fmt.Println("Panic Generated:", any)
	panic(any)
}

func toValueOrPanic(any interface{}) otto.Value {
	ans, err := VM.ToValue(any)
	panicOnError(err)
	return ans
}

func NRGBA2hex(c color.NRGBA) string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

func hex2uint8(hex string) uint8 {
	ans, err := strconv.ParseInt("0x"+hex, 0, 64)
	panicOnError(err)
	return uint8(ans)
}

func hex2NRGBA(hex string) color.NRGBA {
	matched, err := regexp.MatchString("^(#|)[A-Fa-f0-9]{1,8}$", "#FFF")
	panicOnError(err)
	if !matched {
		panic("Invalid color: " + hex)
	}

	ans := color.NRGBA{}
	if hex[0] == '#' {
		hex = hex[1:]
	}
	ans.A = 0xFF
	// Multiplying by 17 "doubles" the algarism. Ex: F => FF
	if len(hex) == 3 || len(hex) == 4 {
		ans.R = 17 * hex2uint8(string(hex[0]))
		ans.G = 17 * hex2uint8(string(hex[1]))
		ans.B = 17 * hex2uint8(string(hex[2]))
		if len(hex) == 4 {
			ans.A = 17 * hex2uint8(string(hex[3]))
		}
	}
	if len(hex) == 6 || len(hex) == 8 {
		ans.R = hex2uint8(hex[0:1])
		ans.G = hex2uint8(hex[2:3])
		ans.B = hex2uint8(hex[4:5])
		if len(hex) == 8 {
			ans.A = hex2uint8(hex[6:7])
		}
	}
	return ans
}

func PreLoad() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	if sdlImg.Init(sdlImg.INIT_JPG|sdlImg.INIT_PNG|sdlImg.INIT_TIF|sdlImg.INIT_WEBP) == 0 {
		panic("failed to init sdl2-img")
	}
}

func Quit() {
	sdl.Quit()
}

func hex2uint32(hex string) uint32 {
	return color2uint32(hex2NRGBA(hex))
}

func color2uint32(c color.Color) uint32 {
	ans := sdl.Color{}
	r, g, b, a := c.RGBA()
	ans.R = uint8(255 * r / 0xffff)
	ans.G = uint8(255 * g / 0xffff)
	ans.B = uint8(255 * b / 0xffff)
	ans.A = uint8(255 * a / 0xffff)

	return ans.Uint32()
}

func RectFitAndCenter(src, dst sdl.Rect) *sdl.Rect {
	max_width := sdl.Rect{
		X: 0,
		Y: 0,
		W: dst.W,
		H: int32(float64(src.H) * (float64(dst.W) / float64(src.W))),
	}
	max_width.X = (dst.W - max_width.W) / 2
	max_width.Y = (dst.H - max_width.H) / 2
	max_height := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(float64(src.W) * (float64(dst.H) / float64(src.H))),
		H: dst.H,
	}
	max_height.X = (dst.W - max_height.W) / 2
	max_height.Y = (dst.H - max_height.H) / 2
	if max_width.H < dst.H {
		return &max_width
	}
	return &max_height
}

func Surface2Rect(surf *sdl.Surface) *sdl.Rect {
	return &sdl.Rect{
		W: surf.W,
		H: surf.H,
		X: 0,
		Y: 0,
	}
}

func RectScale(r *sdl.Rect, scale float64) *sdl.Rect {
	r.W = int32(float64(r.W) * scale)
	r.H = int32(float64(r.H) * scale)
	return r
}

func UserPos2SDLPos(out_rect *sdl.Rect, final_surf *sdl.Surface) *sdl.Rect {
	final_rect := Surface2Rect(final_surf)

	fr_w := float64(final_rect.W) / 2
	fr_h := float64(final_rect.H) / 2
	or_w := float64(out_rect.W) / 2
	or_h := float64(out_rect.H) / 2

	ans_x := float64(out_rect.X) + fr_w - or_w
	ans_y := float64(out_rect.Y) + fr_h - or_h

	out_rect.X = int32(math.Round(ans_x))
	out_rect.Y = int32(math.Round(ans_y))

	return out_rect
}

func RectFitAndCenterInSurf(src sdl.Rect, surf *sdl.Surface) *sdl.Rect {
	return RectFitAndCenter(src, sdl.Rect{0, 0, surf.W, surf.H})
}

func unixMillis() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

func num2float64(num interface{}) float64 {
	switch val := num.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	default:
		panic(fmt.Sprintf("Unable to convert %+v (%+v) to float64", num, reflect.TypeOf(num)))
	}
}

func limit_byte(val int) byte {
	if val < 0 {
		val = 0
	}
	if val > 255 {
		val = 255
	}
	return byte(val)
}
