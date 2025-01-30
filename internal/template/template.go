package template

import (
	"errors"
	"path"
	"strings"
	gt "text/template"
)

const defaultRenderPath = "./temp"

var (
	ErrTemplateNameEmpty = errors.New("template name can not be empty")
	ErrTemplateNotFound  = errors.New("template not found")
)

func Render(name string, data any) (string, error) {
	if name == "" {
		return name, ErrTemplateNameEmpty
	}

	f := path.Join(defaultRenderPath, name)
	t, err := gt.ParseFiles(f)
	if err != nil {
		return "", err
	}

	str := strings.Builder{}
	err = t.Execute(&str, data)
	return str.String(), err
}
