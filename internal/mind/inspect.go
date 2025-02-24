package mind

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Inspect interface {
	Inspect() (io.ReadSeeker, error)
}

type MockScreenInspector struct {
	data  []io.ReadSeeker
	index uint
}

func (m MockScreenInspector) Inspect() (io.ReadSeeker, error) {
	if m.index >= uint(len(m.data)) {
		return nil, errors.New("index out of bound")
	}
	defer func() {
		m.index += 1
	}()
	return m.data[m.index], nil
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

	return MockScreenInspector{
		data:  readers,
		index: 0,
	}
}
