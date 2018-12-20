package AnimaKit

import colorful "github.com/lucasb-eyer/go-colorful"

type ColorMixerSegment struct {
	StartTime float64
	EndTime   float64
	StartVal  colorful.Color
	EndVal    colorful.Color
}

func (self ColorMixerSegment) ValAt(at float64) colorful.Color {
	delta_t := self.EndTime - self.StartTime
	at = at - self.StartTime
	t := at / delta_t
	return self.StartVal.BlendLuv(self.EndVal, t)
}

type ColorMixer struct {
	Segs []ColorMixerSegment
}

func (self ColorMixer) ValAt(at float64) colorful.Color {
	if len(self.Segs) == 0 {
		return colorful.Color{0, 0, 0}
	}

	// Find correct segment
	for _, seg := range self.Segs {
		if seg.StartTime <= at && at < seg.EndTime {
			return seg.ValAt(at)
		}
	}

	// If there is no segment, use the last value as a fixed thing
	return self.Segs[len(self.Segs)-1].EndVal
}
