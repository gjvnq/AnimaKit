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
				fmt.Println("----------------------")
				fmt.Println(sdl.GetKeyName(event.Keysym.Sym))
				fmt.Println("mod", event.Keysym.Mod, event.Keysym.Mod&sdl.KMOD_SHIFT)
				fmt.Println("speed/frame", speed, current_frame)
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_LEFT {
					need_to_redraw = true
					speed = 0
					// [SHIFT] + [←] = first frame (and pause)
					if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
						fmt.Println("First frame")
						current_frame = 0
					}
					// [←] = previous frame
					if event.Keysym.Mod == sdl.KMOD_NONE {
						fmt.Println("Previous frame")
						current_frame--
					}
				}
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_RIGHT {
					need_to_redraw = true
					speed = 0
					// [SHIFT] + [→] = last frame (and pause)
					if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
						fmt.Println("Last frame")
						fmt.Println(TheAnimation.Length, TheAnimation.FPS)
						current_frame = int(TheAnimation.Length * TheAnimation.FPS)
					}
					// [→] = next frame
					if event.Keysym.Mod == sdl.KMOD_NONE {
						fmt.Println("Next frame")
						current_frame++
					}
				}
				// [SPACE] pauses and unpauses the animation
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_SPACE {
					if event.Keysym.Mod == sdl.KMOD_NONE {
						if speed > 0 {
							speed = 0
							fmt.Println("stop playing")
						} else {
							speed = 1
							fmt.Println("fwd")
						}
						playing_last_update = unixMillis()
						// Shift makes it play in reverse
					} else if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
						if speed < 0 {
							speed = 0
							fmt.Println("stop playing")
						} else {
							speed = -1
							fmt.Println("rev")
						}
					}
				}
				fmt.Println(speed, current_frame)
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
			fmt.Println(current_frame)
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
