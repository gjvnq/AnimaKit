package AnimaKit

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

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
	matched, err := regexp.MatchString("^[A-Fa-f0-9]{1,8}$", "#FFF")
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
	ans.R = uint8(r)
	ans.G = uint8(g)
	ans.B = uint8(b)
	ans.A = uint8(a)

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

func RectFitAndCenterInSurf(src sdl.Rect, surf *sdl.Surface) *sdl.Rect {
	return RectFitAndCenter(src, sdl.Rect{0, 0, surf.W, surf.H})
}

func unixMillis() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Millisecond)
}
