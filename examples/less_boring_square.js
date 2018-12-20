Animation.width = 800 // in pixels
Animation.height = 600 // in pixels
Animation.fps = 15 // in frames per second
Animation.length = 60/Animation.fps
Animation.stage = new HiBitStage(64, 64)
Animation.stage.bg = {
                        0: "#FFF",
                        1: "#F00",
                        2: "#0F0",
                        3: "#00F",
                        4: "#F50",
                        30: "#F00"
                      }
println(Animation.stage.bg)