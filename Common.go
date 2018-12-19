package AnimaKit

import (
	"fmt"
	"image/color"
	"strconv"

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
	ans := color.NRGBA{}
	if hex[0] == '#' {
		hex = hex[1:]
	}
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
