package mind_test

import (
	"testing"

	"github.com/farhoud/confidant/internal/mind"
	"github.com/stretchr/testify/assert"
)

func setup() mind.Vision {
	return mind.NewVision()
}

func TestBlindVision(t *testing.T) {
	v := setup()
	box, err := v.Detect("icon", nil)

	assert.Error(t, err)
	assert.Equal(t, mind.ErrBlindVision, err)
	assert.Empty(t, box)
}
