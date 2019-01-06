package AnimaKit

import (
	"github.com/robertkrimen/otto"
	"github.com/veandco/go-sdl2/sdl"
)

type Viz interface {
	DrawOn(frame float64, surface *sdl.Surface) error
}

func get_Viz(id otto.Value) Viz {
	id_int, err := id.ToInteger()
	panicOnError(err)

	return mapperGet(int(id_int)).(Viz)
}
