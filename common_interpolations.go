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

func (self *PositionableBase) pos_parse(sorted_keys []int, key_frame_spec map[int]map[string]interface{}) {
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
}

func (self *ScalableBase) scale_parse(sorted_keys []int, key_frame_spec map[int]map[string]interface{}) {
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
}

func (self *VisibleBase) visible_parse(sorted_keys []int, key_frame_spec map[int]map[string]interface{}) {
	self.Visible.Clear()
	// Ensure default value
	has_defualt := false
	if sorted_keys[0] == 0 {
		_, has_defualt = key_frame_spec[0]["scale"]
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
		self.Visible.Append(InterpSegᐸboolᐳ{
			Frame: float64(key),
			Val:   params["visible"].(bool),
		})
	}
}
