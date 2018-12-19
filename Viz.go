package AnimaKit

import "github.com/veandco/go-sdl2/sdl"

type Viz interface {
	DrawOn(frame int, surface *sdl.Surface) error
}
