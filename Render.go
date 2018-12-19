package AnimaKit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/veandco/go-sdl2/img"
)

func RenderTo(output_path string) {
	output_path, _ = filepath.Abs(output_path)
	fmt.Println("Rendering to directory:", output_path)
	err := os.MkdirAll(output_path, 0755)
	panicOnError(err)

	fmt.Println("Total frames to render:", TheAnimation.Frames)
	fmt.Println("FPS:", TheAnimation.FPS)
	fmt.Println("Length:", TheAnimation.Length)

	surface := TheAnimation.NewSurface()
	fmt.Printf("%+v\n", surface.Format)
	for i := 0; i < TheAnimation.Frames; i++ {
		start_time := unixMillis()

		fmt.Println("Rendering frame:", i)
		frame_filename := fmt.Sprintf("%s/%05d.png", output_path, i)
		TheAnimation.DrawOn(i, surface)
		err := img.SavePNG(surface, frame_filename)
		fmt.Printf("Done frame %d in %0.3f miliceonds\n", i, unixMillis()-start_time)
		panicOnError(err)
	}
}
