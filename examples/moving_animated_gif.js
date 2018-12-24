Animation.width = 800 // in pixels
Animation.height = 600 // in pixels
Animation.fps = 15 // in frames per second
Animation.length = 90/Animation.fps
Animation.stage = new HiBitStage(64, 64)
Animation.stage.bg = {
                        0: "#89CFF0",
                        90: "#CCCCFF"
                      }
println(Animation.stage.bg)

var gif = new Image("bandits.gif")
Animation.stage.place(gif, {
  0: {x: 0, y: 0, scale: 1},
  90: {x 128, y: 128, scale: 0.5}
})
