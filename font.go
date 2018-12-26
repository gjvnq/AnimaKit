package AnimaKit

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func loadInternalFont(filename string, size int) (*ttf.Font, error) {
	data, err := Asset("res/" + filename)
	if err != nil {
		panic(err)
	}
	rwops, err := sdl.RWFromMem(data)
	if err != nil {
		panic(err)
	}
	return ttf.OpenFontRW(rwops, 0, size)
}
