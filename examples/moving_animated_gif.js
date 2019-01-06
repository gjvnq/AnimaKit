Animation.width = 800 // in pixels
Animation.height = 600 // in pixels
Animation.fps = 15 // in frames per second
Animation.length = 90/Animation.fps
Animation.stage = new HiBitStage(512, 512)
Animation.stage.bg = {
                        0: "#89CFF0",
                        90: "#CCCCFF"
                      }
// println(Animation.stage.bg)

var gif = new GIF("bandits.gif")
gif.keyframes = {
  1: {visible: false},
  5: {x: 0, y: 0, scale: 1, visible: true},
  90: {x: 128, y: 128, scale: 0.5}
}
Animation.stage.place(gif)
