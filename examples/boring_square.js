Animation.width = 1920 // in pixels
Animation.height = 1080 // in pixels
Animation.fps = 1080 // in frames per second
println(Animation.stage)
Animation.stage = new HiBitStage(256, 240)
println(Animation.stage)
println(Animation.stage.bg)
Animation.stage.bg = "#00F"
println([Animation.stage.bg, Animation, Animation.stage.bg])