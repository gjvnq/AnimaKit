Animation.width = 800 // in pixels
Animation.height = 600 // in pixels
Animation.fps = 30 // in frames per second
Animation.length = 10
Animation.stage = new HiBitStage(64, 64)
Animation.stage.bg = {
                        0: "#00F",
                        150: "#0F0",
                        300: "#F00",
                      }
println(Animation.stage.bg)