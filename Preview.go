package AnimaKit

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func PreviewWindow() {
	window, err := sdl.CreateWindow("AnimaKit", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	current_frame := 0
	need_to_redraw := true
	running := true
	for running {
		for event_raw := sdl.PollEvent(); event_raw != nil; event_raw = sdl.PollEvent() {
			switch event := event_raw.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.WindowEvent:
				if event.Type != sdl.WINDOWEVENT_MOVED {
					need_to_redraw = true
				}
			}
		}
		if need_to_redraw {
			surface, err := window.GetSurface()
			if err != nil {
				panic(err)
			}

			need_to_redraw = false
			fmt.Println(surface)
			fmt.Printf("%+v\n", surface.Format)
			TheAnimation.DrawOn(current_frame, surface)
			window.UpdateSurface()
			// img.SavePNG(surface, "last.png")
		}
	}
}
