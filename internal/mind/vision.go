package mind

import (
	"errors"
	"io"
)

var ErrBlindVision = errors.New("input is null we are blind")

type Box struct {
	TRPoint [2]uint
	BLPoint [2]uint
}

type Vision interface {
	Detect(description string, input io.Reader) (Box, error)
}

type vision struct{}

func (v vision) Detect(d string, i io.Reader) (Box, error) {
	if i == nil {
		return Box{}, ErrBlindVision
	}
	return Box{}, nil
}

func NewVision() Vision {
	return &vision{}
}
