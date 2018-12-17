package AnimaKit

type KeySegment interface {
	StartVal() float64
	StartTime() float64
	EndTime() float64
	EndVal() float64
	Set(start_time, end_time, start_val, end_val float64)
	ValAt(at float64) float64
}

type Interpolator struct {
	Segs []KeySegment
}

func NewInterpolator() *Interpolator {
	ans := new(Interpolator)
	ans.Segs = make([]KeySegment, 0)
	return ans
}

func (self *Interpolator) ValAt(at float64) float64 {
	if len(self.Segs) == 0 {
		return 0
	}

	// Find correct segment
	for _, seg := range self.Segs {
		if seg.StartTime() <= at && at < seg.EndTime() {
			return seg.ValAt(at)
		}
	}

	// If there is no segment, use the last value as a fixed thing
	return self.Segs[len(self.Segs)-1].EndVal()
}

type BaseSegment struct {
	start_time float64
	start_val  float64
	end_time   float64
	end_val    float64
}

func (self BaseSegment) StartTime() float64 {
	return self.start_time
}

func (self BaseSegment) StartVal() float64 {
	return self.start_val
}

func (self BaseSegment) EndVal() float64 {
	return self.end_val
}

func (self BaseSegment) EndTime() float64 {
	return self.end_time
}

type LinearSegment struct {
	BaseSegment
	m float64
	b float64
}

func NewLinearSegment(start_time, end_time, start_val, end_val float64) *LinearSegment {
	ans := new(LinearSegment)
	ans.Set(start_time, end_time, start_val, end_val)
	return ans
}

func (self *LinearSegment) Set(start_time, end_time, start_val, end_val float64) {
	self.start_time = start_time
	self.end_time = end_time
	self.start_val = start_val
	self.end_val = end_val

	end_time = end_time - start_time
	start_time = 0

	self.m = (start_val - end_val) / (start_time - end_time)
	self.b = start_val + self.m*start_time
}

func (self LinearSegment) ValAt(time float64) float64 {
	time = time - self.start_time
	return time*self.m + self.b
}
