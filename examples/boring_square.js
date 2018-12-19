Animation.width = 800 // in pixels
Animation.height = 600 // in pixels
Animation.fps = 30 // in frames per second
Animation.length = 1
println(Animation.fps, Animation.stage)
Animation.stage = new HiBitStage(64, 64)
println(Animation.stage)
println(Animation.stage.bg)
Animation.stage.bg = "#00F"
println([Animation.stage.bg, Animation, Animation.stage.bg])