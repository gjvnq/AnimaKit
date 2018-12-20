package AnimaKit

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/veandco/go-sdl2/img"
)

var timeProgramStart float64

func init() {
	timeProgramStart = unixMillis()
}

func RenderTo(output_path string, n_workers int) {
	output_path, _ = filepath.Abs(output_path)
	fmt.Println("Rendering to directory:", output_path)
	err := os.MkdirAll(output_path, 0755)
	panicOnError(err)

	fmt.Println("Total frames to render:", TheAnimation.Frames)
	fmt.Println("FPS:", TheAnimation.FPS)
	fmt.Println("Length:", TheAnimation.Length)
	fmt.Println("Number of workers:", n_workers)

	wg := new(sync.WaitGroup)
	wg.Add(n_workers)
	var frames_to_do = make(chan int, 64)
	for i := 0; i < n_workers; i++ {
		go renderWorker(output_path, frames_to_do, wg)
	}
	for i := 0; i < TheAnimation.Frames; i++ {
		frames_to_do <- i
	}
	// Tell our workers that all frames have been requested
	for i := 0; i < n_workers; i++ {
		frames_to_do <- -1
	}
	full_time := unixMillis() - timeProgramStart
	wg.Wait()
	fmt.Printf("Full render finished in %f milliseconds. Average time per frame: %7.3f \n", full_time, full_time/float64(TheAnimation.Frames))
}

func renderWorker(dir string, input chan int, wg *sync.WaitGroup) {
	surface := TheAnimation.NewSurface()
	fmt.Println("Started worker")
	for {
		frame := <-input
		if frame == -1 {
			break
		}
		start_time := unixMillis()

		filename := fmt.Sprintf("%s/%05d.png", dir, frame)
		TheAnimation.DrawOn(frame, surface)
		err := img.SavePNG(surface, filename)
		fmt.Printf("Done frame %5d in %7.3f milliseconds\n", frame, unixMillis()-start_time)
		panicOnError(err)
	}
	fmt.Println("Finished worker")
	wg.Done()
}
