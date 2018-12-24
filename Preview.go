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
	speed := 0
	playing_last_update := unixMillis()
	for running {
		for event_raw := sdl.PollEvent(); event_raw != nil; event_raw = sdl.PollEvent() {
			switch event := event_raw.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				// Events using the [←] key
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_LEFT {
					need_to_redraw = true
					speed = 0
					// [SHIFT] + [←] = first frame (and pause)
					if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
						current_frame = 0
					}
					// [←] = previous frame
					if event.Keysym.Mod == sdl.KMOD_NONE {
						current_frame--
					}
				}
				// Events using the [→] key
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_RIGHT {
					need_to_redraw = true
					speed = 0
					// [SHIFT] + [→] = last frame (and pause)
					if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
						current_frame = int(TheAnimation.Length * TheAnimation.FPS)
					}
					// [→] = next frame
					if event.Keysym.Mod == sdl.KMOD_NONE {
						current_frame++
					}
				}
				// [SPACE] pauses and unpauses the animation ([SHIFT] reverses it)
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_SPACE {
					if event.Keysym.Mod == sdl.KMOD_NONE {
						if speed > 0 {
							speed = 0
						} else {
							speed = 1
						}
						playing_last_update = unixMillis()
					} else if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
						if speed < 0 {
							speed = 0
						} else {
							speed = -1
						}
					}
				}
			}
		}
		if speed != 0 && unixMillis()-playing_last_update >= 1000/TheAnimation.FPS {
			current_frame += speed
			need_to_redraw = true
			playing_last_update = unixMillis()
		}
		if need_to_redraw {
			if current_frame < 0 {
				current_frame = 0
				speed = 0
			}
			fmt.Println("Frame = ", current_frame)
			surface, err := window.GetSurface()
			if err != nil {
				panic(err)
			}

			need_to_redraw = false
			TheAnimation.DrawOn(current_frame, surface)
			window.UpdateSurface()
		}
	}
}
