package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gjvnq/AnimaKit"
	"github.com/integrii/flaggy"
)

var ArgInput string = ""
var ArgOutputDir string = ""
var ArgOutputFile string = ""
var ArgFfmpeg string = "ffmpeg"
var ArgNWorkers int = 8

func init() {
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("AnimaKit")
	flaggy.SetVersion(AnimaKit.VERSION)
	flaggy.SetDescription("A simple animation renderer")
	flaggy.DefaultParser.ShowHelpOnUnexpected = false
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/gjvnq/AnimaKit"
	flaggy.Int(&ArgNWorkers, "n", "n-workers", "Number of parallel rendering workers")
	flaggy.String(&ArgFfmpeg, "", "ffmpeg", "Path to ffmpeg CLI tool")
	flaggy.AddPositionalValue(&ArgInput, "input", 1, true, "Animation script file (.js)")
	flaggy.AddPositionalValue(&ArgOutputDir, "frames-dir", 2, false, "Directory to store the output video frames")
	flaggy.AddPositionalValue(&ArgOutputFile, "video-output", 3, false, "Output video file (requires ffmpeg CLI tool)")
	flaggy.Parse()
}

func main() {
	// Load SDL and stuff
	AnimaKit.PreLoad()
	defer AnimaKit.Quit()

	// Load animation
	_, err := AnimaKit.LoadScriptFromFile(ArgInput)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// fmt.Println("Use [←] and [→] to move frame by frame and [SPACE] to play")
	if ArgOutputDir == "" {
		AnimaKit.PreviewWindow()
	} else {
		// Render frames
		AnimaKit.RenderTo(ArgOutputDir, ArgNWorkers)

		// Use ffmpeg to convert it to a video file if the user requested such a thing
		if ArgOutputFile != "" {
			cmd := exec.Command(
				ArgFfmpeg,
				"-y", // overwrite output file
				"-framerate",
				fmt.Sprintf("%f", AnimaKit.TheAnimation.FPS),
				"-i",
				ArgOutputDir+"/%05d.png",
				ArgOutputFile)
			// Force color output
			cmd.Env = append(os.Environ(),
				"AV_LOG_FORCE_COLOR=true",
			)
			// Ensure output will be visible
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			fmt.Println("Running:", cmd.Args)
			fmt.Println("--------------------------------------------------------------------------------")
			cmd.Run()
		}
	}
}
