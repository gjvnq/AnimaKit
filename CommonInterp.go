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
	self.X.Parse("x", 0, LinearInterpolatorᐸF64ᐳ, sorted_keys, key_frame_spec)
	self.Y.Parse("y", 0, LinearInterpolatorᐸF64ᐳ, sorted_keys, key_frame_spec)
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
	self.Scale.Parse("scale", 1, LinearInterpolatorᐸF64ᐳ, sorted_keys, key_frame_spec)
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
