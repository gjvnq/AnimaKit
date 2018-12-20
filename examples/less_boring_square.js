Animation.width = 800 // in pixels
Animation.height = 600 // in pixels
Animation.fps = 15 // in frames per second
Animation.length = 90/Animation.fps
Animation.stage = new HiBitStage(64, 64)
Animation.stage.bg = {
                        0: "#F00",
                        30: "#0F0",
                        60: "#00F"
                      }
println(Animation.stage.bg)