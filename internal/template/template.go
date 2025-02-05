package template

import (
	"errors"
	"path"
	"strings"
	gt "text/template"
)

var (
	ErrTemplateNameEmpty = errors.New("template name can not be empty")
	ErrTemplateNotFound  = errors.New("no such file or directory")
)

type Template interface {
	Render(string, any) (string, error)
}

type tmpl struct {
	path string
}

func (t tmpl) Render(name string, data any) (string, error) {
	if name == "" {
		return name, ErrTemplateNameEmpty
	}

	f := path.Join(t.path, name)
	tmpl, err := gt.ParseFiles(f)
	if err != nil {
		err = errors.Unwrap(err)
		return "", err
	}

	str := strings.Builder{}
	err = tmpl.Execute(&str, data)
	return str.String(), err
}

func NewTemplateEngine(path string) Template {
	return &tmpl{
		path: path,
	}
}
