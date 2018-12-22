package AnimaKit

import (
	"image/color"
	"sort"
	"strconv"

	"github.com/robertkrimen/otto"
	colorful "gopkg.in/lucasb-eyer/go-colorful.v1"
)

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
	return self.StartVal.BlendHcl(self.EndVal, t).Clamped()
}

type ColorMixer struct {
	Segs []ColorMixerSegment
}

func (self ColorMixer) ValAt(at float64) color.Color {
	if len(self.Segs) == 0 {
		return color.RGBA{0, 0, 0, 0}
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

func (self *ColorMixer) Clear() {
	self.Segs = nil
}

func (self *ColorMixer) SetFixed(hex string) {
	self.Segs = make([]ColorMixerSegment, 1)
	c, err := colorful.Hex(hex)
	panicOnError(err)
	self.Segs[0].StartVal = c
	self.Segs[0].EndVal = c
}

func (self *ColorMixer) FromValue(value otto.Value) {
	if value.IsString() {
		hex, err := value.ToString()
		panicOnError(err)
		self.SetFixed(hex)
	} else {
		self.Clear()

		// Get and sort keys
		keys := make([]int, 0)
		for _, key := range value.Object().Keys() {
			val, err := strconv.Atoi(key)
			panicOnError(err)
			keys = append(keys, val)
		}
		sort.Ints(keys)

		// Pre add first segment
		if keys[0] != 0 {
			self.Segs = append(self.Segs, ColorMixerSegment{
				StartVal:  colorful.Color{0, 0, 0},
				StartTime: 0,
			})
		}

		// Get values for each key frame
		tmp, err := value.Export()
		panicOnError(err)
		obj := tmp.(map[string]interface{})
		// Go in order
		for _, key_int := range keys {
			key := strconv.Itoa(key_int)
			i := len(self.Segs)

			// Convert value to colour
			c, err := colorful.Hex(obj[key].(string))
			panicOnError(err)

			key_f64, err := strconv.ParseFloat(key, 64)
			panicOnError(err)

			self.Segs = append(self.Segs, ColorMixerSegment{
				StartVal:  c,
				EndVal:    c,
				StartTime: key_f64,
				EndTime:   float64(TheAnimation.Frames),
			})
			if i != 0 {
				self.Segs[i-1].EndVal = c
				self.Segs[i-1].EndTime = key_f64
			}
		}
	}
}
