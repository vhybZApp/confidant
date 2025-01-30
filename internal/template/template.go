package template

import (
	"errors"
	"path"
	"strings"
	gt "text/template"
)

const defaultRenderPath = "./tmpl"

var (
	ErrTemplateNameEmpty = errors.New("template name can not be empty")
	ErrTemplateNotFound  = errors.New("no such file or directory")
)

func Render(name string, data any) (string, error) {
	if name == "" {
		return name, ErrTemplateNameEmpty
	}

	f := path.Join(defaultRenderPath, name)
	t, err := gt.ParseFiles(f)
	if err != nil {
		err = errors.Unwrap(err)
		return "", err
	}

	str := strings.Builder{}
	err = t.Execute(&str, data)
	return str.String(), err
}
