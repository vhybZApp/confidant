package mind

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/go-vgo/robotgo"
)

type Inspect interface {
	Inspect() (io.ReadSeeker, error)
}

type MockScreenInspector struct {
	data  []io.ReadSeeker
	index uint
}

func (m *MockScreenInspector) Inspect() (io.ReadSeeker, error) {
	if m.index >= uint(len(m.data)) {
		return nil, errors.New("index out of bound")
	}
	data := m.data[m.index]

	m.index += 1
	return data, nil
}

func NewMockScreenInspector(paths []string) Inspect {
	readers := make([]io.ReadSeeker, len(paths))

	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			panic(fmt.Sprintf("failed to open file: %s", path))
		}
		readers[i] = file
	}

	return &MockScreenInspector{
		data:  readers,
		index: 0,
	}
}

type RobotScreenInspector struct{}

func (r RobotScreenInspector) Inspect() (io.ReadSeeker, error) {
	robotgo.DisplayID = 0
	img, err := robotgo.Capture()
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}

	return bytes.NewReader(buff.Bytes()), nil
}

func NewRobotScreenInspector() Inspect {
	return RobotScreenInspector{}
}
