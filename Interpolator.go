package AnimaKit

type Interpolator func(start_time, end_time, start_val, end_val, at_time float64) float64

func LinearInterpolator(start_time, end_time, start_val, end_val, at_time float64) float64 {
	end_time = end_time - start_time
	start_time = 0

	m := (start_val - end_val) / (start_time - end_time)
	b := start_val + m*start_time

	at_time = at_time - start_time
	return at_time*m + b
}

type InterpSeg struct {
	StartVal     float64
	EndVal       float64
	StartTime    float64
	EndTime      float64
	Interpolator Interpolator
}

func NewInterpSeg(start_time, end_time, start_val, end_val float64, interpolator Interpolator) InterpSeg {
	ans := InterpSeg{}
	ans.StartTime = start_time
	ans.StartTime = end_time
	ans.StartVal = start_val
	ans.EndVal = end_val
	ans.Interpolator = interpolator
	if interpolator == nil {
		ans.Interpolator = LinearInterpolator
	}

	return ans
}

func (self InterpSeg) ValAtᐸintᐳ(at_time int) float64 {
	return self.ValAt(float64(at_time))
}

func (self InterpSeg) ValAt(at_time float64) float64 {
	f := self.Interpolator
	if f == nil {
		f = LinearInterpolator
	}
	return f(
		self.StartTime,
		self.EndTime,
		self.StartVal,
		self.EndVal,
		at_time)
}
