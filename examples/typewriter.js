Animation.width = 1920 // in pixels
Animation.height = 1080 // in pixels
Animation.length = 15 // in pixels

// HiBitStage(width, height, framerate)
Animation.stage = new HiBitStage(256, 240, 30)

// Set the stage to transparent with pink fallback
Animation.stage.bg = "#F0FF"

var typewriter = new TypeWriter({
  msg: "Hello World",
  font: "Unifont",
  fg: "#FFF"})

Animation.stage.place(typewriter, [
  {frame: 0, visible: false, x: 0, y: 0}, // coordinates origin in the center of the frame
  {frame: 30, cur: 0, speed: 30}, // visible = true unless otherwise; cur = cursor position; speed = characters per frame
  {cur: 6, speed: 15} // Speed up after the first word
])
