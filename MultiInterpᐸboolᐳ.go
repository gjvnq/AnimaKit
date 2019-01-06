package AnimaKit

type InterpSegᐸboolᐳ struct {
	Frame float64
	Val   bool
}

func NewInterpSegᐸboolᐳ(frame float64, val bool) InterpSegᐸboolᐳ {
	return InterpSegᐸboolᐳ{
		Frame: frame,
		Val:   val,
	}
}

type MultiInterpᐸboolᐳ struct {
	Segs []InterpSegᐸboolᐳ
}

func (self MultiInterpᐸboolᐳ) ValAt(frame float64) bool {
	// False by default
	if len(self.Segs) == 0 {
		return false
	}

	i := 0
	for ; i < len(self.Segs)-1; i++ {
		this := self.Segs[i]
		next := self.Segs[i+1]
		// Is this the correct segment?
		if frame < next.Frame {
			return this.Val
		}
	}
	// Returns the last segment as it is the best option avaliable
	return self.Segs[len(self.Segs)-1].Val
}

func (self *MultiInterpᐸboolᐳ) Clear() {
	self.Segs = make([]InterpSegᐸboolᐳ, 0)
}

func (self *MultiInterpᐸboolᐳ) Append(new_seg InterpSegᐸboolᐳ) {
	self.Segs = append(self.Segs, new_seg)
}
