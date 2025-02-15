package mind

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Inspect interface {
	Inspect() (io.Reader, error)
}

type MockScreenInspector struct {
	data  []io.Reader
	index uint
}

func (m MockScreenInspector) Inspect() (io.Reader, error) {
	if m.index >= uint(len(m.data)) {
		return nil, errors.New("index out of bound")
	}
	defer func() {
		m.index += 1
	}()
	return m.data[m.index], nil
}

func NewMockScreenInspector(paths []string) Inspect {
	var readers []io.Reader

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			panic(fmt.Sprintf("failed to open file: %s", path))
		}
		readers = append(readers, file)
	}

	return MockScreenInspector{
		data:  readers,
		index: 0,
	}
}
