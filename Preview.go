package AnimaKit

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

func PreviewWindow() {
	window, err := sdl.CreateWindow("AnimaKit", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		600, 700, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	current_frame := 0
	need_to_redraw := true
	running := true
	speed := 0
	playing_last_update := unixMillis()
	frame_goto_buf := ""
	frame_goto_flag := false
	font, err := loadInternalFont("unifont.ttf", 16)
	if err != nil {
		panic(err)
	}
	for running {
		for event_raw := sdl.PollEvent(); event_raw != nil; event_raw = sdl.PollEvent() {
			switch event := event_raw.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.WindowEvent:
				need_to_redraw = true
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
					if event.Keysym.Mod&sdl.KMOD_SHIFT == 0 {
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
					if event.Keysym.Mod&sdl.KMOD_SHIFT == 0 {
						current_frame++
					}
				}
				// [SPACE] pauses and unpauses the animation ([SHIFT] reverses it)
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_SPACE {
					// No [SHIFT]
					if event.Keysym.Mod&sdl.KMOD_SHIFT == 0 {
						if speed > 0 {
							speed = 0
						} else {
							speed = 1
						}
					} else if event.Keysym.Mod&sdl.KMOD_SHIFT != 0 { // with [SHIFT]
						if speed < 0 {
							speed = 0
						} else {
							speed = -1
						}
					}
					need_to_redraw = true
					playing_last_update = unixMillis()
				}

				// [J] [0-9] [ENTER] allow to jump to frame
				// Enable frame typing
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_j {
					frame_goto_flag = true
					frame_goto_buf = ""
				}

				// [J] [0-9] [ENTER] allow to jump to frame
				// Append digit
				if event.Type == sdl.KEYUP && frame_goto_flag {
					switch event.Keysym.Sym {
					case sdl.K_0:
						frame_goto_buf += "0"
					case sdl.K_1:
						frame_goto_buf += "1"
					case sdl.K_2:
						frame_goto_buf += "2"
					case sdl.K_3:
						frame_goto_buf += "3"
					case sdl.K_4:
						frame_goto_buf += "4"
					case sdl.K_5:
						frame_goto_buf += "5"
					case sdl.K_6:
						frame_goto_buf += "6"
					case sdl.K_7:
						frame_goto_buf += "7"
					case sdl.K_8:
						frame_goto_buf += "8"
					case sdl.K_9:
						frame_goto_buf += "9"
					default:
					}
				}

				// [J] [0-9] [ENTER] allow to jump to frame
				// Actually go to said frame
				if event.Type == sdl.KEYUP && event.Keysym.Sym == sdl.K_RETURN {
					new_frame, err := strconv.Atoi(frame_goto_buf)
					if err != nil {
						TheLog.ErrorF("Failed to go to understand frame %s", frame_goto_buf)
					} else {
						current_frame = new_frame
						need_to_redraw = true
						speed = 0
					}
					frame_goto_flag = false
					frame_goto_buf = ""
				}
			}
		}
		if speed != 0 && unixMillis()-playing_last_update >= 1000/TheAnimation.FPS {
			current_frame += speed
			need_to_redraw = true
			playing_last_update = unixMillis()
		}
		// Redraw only if necessary
		if need_to_redraw {
			if current_frame < 0 {
				current_frame = 0
				speed = 0
			}
			surface, err := window.GetSurface()
			if err != nil {
				panic(err)
			}

			need_to_redraw = false
			TheAnimation.DrawOn(float64(current_frame), surface)
			// Print debug info
			txt := fmt.Sprintf(
				"Frame: %d Time: %f Speed: %d\nSee shortcuts.md",
				current_frame,
				float64(current_frame)/TheAnimation.FPS,
				speed)
			txt_surf, err := font.RenderUTF8BlendedWrapped(txt, sdl.Color{255, 255, 255, 255}, int(surface.W))
			if err != nil {
				panic(err)
			}
			txt_surf.BlitScaled(nil, surface, &txt_surf.ClipRect)
			// Finish it
			window.UpdateSurface()
		}
	}
}
