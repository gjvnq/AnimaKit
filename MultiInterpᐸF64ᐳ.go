package AnimaKit

import "fmt"

type InterpolatorFuncᐸF64ᐳ func(start_frame, end_frame, start_val, end_val, at_frame float64) float64

func LinearInterpolatorᐸF64ᐳ(start_frame, end_frame, start_val, end_val, at_frame float64) float64 {
	end_frame = end_frame - start_frame

	m := (start_val - end_val) / (0 - end_frame)
	b := start_val + m*0

	at_frame = at_frame - start_frame
	return at_frame*m + b
}

type InterpSegᐸF64ᐳ struct {
	StartVal   float64
	EndVal     float64
	StartFrame float64
	EndFrame   float64
	InterpFunc InterpolatorFuncᐸF64ᐳ
}

func NewInterpSegᐸF64ᐳ(start_frame, end_frame, start_val, end_val float64, interp_func InterpolatorFuncᐸF64ᐳ) InterpSegᐸF64ᐳ {
	ans := InterpSegᐸF64ᐳ{}
	if start_frame == end_frame {
		panic(fmt.Sprintf("start_frame (%f) cannot be equal to end_frame (%f)", start_frame, end_frame))
	}
	ans.StartFrame = start_frame
	ans.StartFrame = end_frame
	ans.StartVal = start_val
	ans.EndVal = end_val
	ans.InterpFunc = interp_func
	if interp_func == nil {
		ans.InterpFunc = LinearInterpolatorᐸF64ᐳ
	}

	return ans
}

func (self InterpSegᐸF64ᐳ) ValAt(at_frame float64) float64 {
	f := self.InterpFunc
	if f == nil {
		f = LinearInterpolatorᐸF64ᐳ
	}

	return f(
		self.StartFrame,
		self.EndFrame,
		self.StartVal,
		self.EndVal,
		at_frame)
}

type MultiInterpᐸF64ᐳ struct {
	Segs []InterpSegᐸF64ᐳ
}

func (self MultiInterpᐸF64ᐳ) ValAt(at_frame float64) float64 {
	// 0 by default
	if len(self.Segs) == 0 {
		return 0
	}

	i := 0
	for ; i < len(self.Segs)-1; i++ {
		this := self.Segs[i]
		next := self.Segs[i+1]
		// Is this the correct segment?
		if at_frame < next.StartFrame {
			return this.ValAt(at_frame)
		}
	}
	// Returns the last segment as it is the best option avaliable
	return self.Segs[len(self.Segs)-1].ValAt(at_frame)
}

func (self *MultiInterpᐸF64ᐳ) Clear() {
	self.Segs = make([]InterpSegᐸF64ᐳ, 0)
}

func (self *MultiInterpᐸF64ᐳ) Append(new_seg InterpSegᐸF64ᐳ) {
	self.Segs = append(self.Segs, new_seg)
}

func (self *MultiInterpᐸF64ᐳ) FixLast(end_frame, end_val float64) {
	if len(self.Segs) == 0 {
		return
	}
	self.Segs[len(self.Segs)-1].EndFrame = end_frame
	self.Segs[len(self.Segs)-1].EndVal = end_val
}
