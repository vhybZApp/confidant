package mind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordination(t *testing.T) {
	ai := AnnotatedImage{
		Annotations: []Annotation{{
			ID:         0,
			BoundedBox: [4]float64{0.726291835308075, 0.9287018775939941, 0.7680469751358032, 1},
		}},
		Width:  1280,
		Height: 715,
	}

	t.Log(ai)
	expected_bb := []int{934, 668, 40, 40}
	x, y := ai.BoundedBox(0)
	t.Logf("x: %d, y: %d", x, y)
	xInBound := x >= expected_bb[0] && x <= expected_bb[0]+expected_bb[2]
	yInBound := y >= expected_bb[1] && y <= expected_bb[1]+expected_bb[3]
	assert.True(t, xInBound && yInBound, "it should be inbound")
}
