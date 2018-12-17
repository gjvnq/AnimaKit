package AnimaKit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterpolator(t *testing.T) {
	interp := NewInterpolator()

	assert.Equal(t, 0.0, interp.ValAt(-42))

	interp.Segs = append(interp.Segs, NewLinearSegment(0, 2, -2, 2))
	interp.Segs = append(interp.Segs, NewLinearSegment(2, 3, 10, 7))

	assert.Equal(t, -2.0, interp.ValAt(0))
	assert.Equal(t, -1.0, interp.ValAt(0.5))
	assert.Equal(t, 0.0, interp.ValAt(1))
	assert.Equal(t, 1.0, interp.ValAt(1.5))
	assert.Equal(t, 1.98, interp.ValAt(1.99))
	assert.Equal(t, 10.0, interp.ValAt(2))
	assert.Equal(t, 7.0, interp.ValAt(3))
	assert.Equal(t, 7.0, interp.ValAt(999))
}
