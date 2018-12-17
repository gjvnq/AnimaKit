function prepare()
	stage = HiBitStage:new(256, 240, 30)
	return stage
end

function animate(stage)
	-- Set the stage to transparent with pink fallback
	stage.bg = "#F0F0"

	typewriter = TypeWriter(
		msg: "Hello World",
		font: "Unifont",
		fg: "#FFF")
end