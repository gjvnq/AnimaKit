package AnimaKit

type PositionableBase struct {
	X MultiInterpᐸF64ᐳ
	Y MultiInterpᐸF64ᐳ
}

type ScalableBase struct {
	Scale MultiInterpᐸF64ᐳ
}

type VisibleBase struct {
	Visible MultiInterpᐸboolᐳ
}

func (self *PositionableBase) pos_init() {
	self.X.Clear()
	self.Y.Clear()
	self.X.Append(InterpSegᐸF64ᐳ{
		StartFrame: 0,
		EndFrame:   PosInf,
		StartVal:   1,
		EndVal:     1,
	})
	self.Y.Append(InterpSegᐸF64ᐳ{
		StartFrame: 0,
		EndFrame:   PosInf,
		StartVal:   1,
		EndVal:     1,
	})
}

func (self *PositionableBase) pos_parse(sorted_keys []float64, key_frame_spec map[float64]map[string]interface{}) {
	TheLog.DebugF("pos_parse %d keys", len(sorted_keys))
	defer TheLog.DebugF("[FINISHED] pos_parse %d keys", len(sorted_keys))

	self.X.Clear()
	self.Y.Clear()
	has_defualt_x := false
	has_defualt_y := false
	if sorted_keys[0] == 0 {
		_, has_defualt_x = key_frame_spec[0]["x"]
		_, has_defualt_y = key_frame_spec[0]["y"]
	}
	if !has_defualt_x {
		has_defualt_x = true
		self.X.Append(InterpSegᐸF64ᐳ{
			StartFrame: 0,
			EndFrame:   PosInf,
			StartVal:   1,
			EndVal:     1,
		})
	}
	if !has_defualt_y {
		has_defualt_y = true
		self.Y.Append(InterpSegᐸF64ᐳ{
			StartFrame: 0,
			EndFrame:   PosInf,
			StartVal:   1,
			EndVal:     1,
		})
	}
	for _, key := range sorted_keys {
		params := key_frame_spec[key]

		if _, ok := params["x"]; ok {
			// Add key frame
			self.X.FixLast(key, num2float64(params["x"]))
			self.X.Append(InterpSegᐸF64ᐳ{
				StartFrame: key,
				EndFrame:   PosInf,
				StartVal:   num2float64(params["x"]),
				EndVal:     num2float64(params["x"]),
			})
			TheLog.DebugF("x = %f at frame %.0f Full params: %#+v", num2float64(params["x"]), key, params)
		}
		if _, ok := params["y"]; ok {
			// Add key frame
			self.Y.FixLast(key, num2float64(params["y"]))
			self.Y.Append(InterpSegᐸF64ᐳ{
				StartFrame: key,
				EndFrame:   PosInf,
				StartVal:   num2float64(params["y"]),
				EndVal:     num2float64(params["y"]),
			})
			TheLog.DebugF("y = %f at frame %.0f Full params: %#+v", num2float64(params["y"]), key, params)
		}
	}
	TheLog.DebugF("PositionableBase.X.Segs: %#+v", self.X.Segs)
	TheLog.DebugF("PositionableBase.Y.Segs: %#+v", self.Y.Segs)
}

func (self *ScalableBase) scale_init() {
	self.Scale.Clear()
	self.Scale.Append(InterpSegᐸF64ᐳ{
		StartFrame: 0,
		EndFrame:   PosInf,
		StartVal:   1,
		EndVal:     1,
	})
}

func (self *ScalableBase) scale_parse(sorted_keys []float64, key_frame_spec map[float64]map[string]interface{}) {
	TheLog.DebugF("scale_parse %d keys", len(sorted_keys))
	defer TheLog.DebugF("[FINISHED] scale_parse %d keys", len(sorted_keys))

	self.Scale.Clear()
	// Ensure default value
	has_defualt := false
	if sorted_keys[0] == 0 {
		_, has_defualt = key_frame_spec[0]["scale"]
	}
	if !has_defualt {
		has_defualt = true
		self.Scale.Append(InterpSegᐸF64ᐳ{
			StartFrame: 0,
			EndFrame:   PosInf,
			StartVal:   1,
			EndVal:     1,
		})
	}
	for _, key := range sorted_keys {
		params := key_frame_spec[key]

		// Ignore key frames that aren't about us
		if _, ok := params["scale"]; !ok {
			continue
		}

		// Add key frame
		self.Scale.FixLast(key, num2float64(params["scale"]))
		self.Scale.Append(InterpSegᐸF64ᐳ{
			StartFrame: key,
			EndFrame:   PosInf,
			StartVal:   num2float64(params["scale"]),
			EndVal:     num2float64(params["scale"]),
		})
		TheLog.DebugF("scale = %f at frame %.0f Full params: %#+v", num2float64(params["scale"]), key, params)
	}
	TheLog.DebugF("ScalableBase.Segs: %#+v", self.Scale.Segs)
}

func (self *VisibleBase) visible_init() {
	self.Visible.Clear()
	self.Visible.Append(InterpSegᐸboolᐳ{
		Frame: 0,
		Val:   true,
	})
}

func (self *VisibleBase) visible_parse(sorted_keys []float64, key_frame_spec map[float64]map[string]interface{}) {
	TheLog.DebugF("visible_parse %d keys", len(sorted_keys))
	defer TheLog.DebugF("[FINISHED] visible_parse %d keys", len(sorted_keys))

	self.Visible.Clear()
	// Ensure default value
	has_defualt := false
	if sorted_keys[0] == 0 {
		_, has_defualt = key_frame_spec[0]["visible"]
	}
	if !has_defualt {
		has_defualt = true
		self.Visible.Append(InterpSegᐸboolᐳ{
			Frame: 0,
			Val:   true,
		})
	}
	for _, key := range sorted_keys {
		params := key_frame_spec[key]

		// Ignore key frames that aren't about us
		if _, ok := params["visible"]; !ok {
			continue
		}

		// Add key frame
		TheLog.DebugF("visible = %t at frame %.0f Full params: %#+v", params["visible"].(bool), key, params)
		self.Visible.Append(InterpSegᐸboolᐳ{
			Frame: key,
			Val:   params["visible"].(bool),
		})
	}
	TheLog.DebugF("VisibleBase.Segs: %#+v", self.Visible.Segs)
}
