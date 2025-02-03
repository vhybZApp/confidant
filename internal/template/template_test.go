package template_test

import (
	"testing"

	"github.com/farhoud/confidant/internal/template"
	"github.com/stretchr/testify/assert"
)

func TestEmptyNameRender(t *testing.T) {
	tmpl := template.NewTemplateEngine("./test_data")
	_, err := tmpl.Render("", nil)

	assert.EqualError(t, err, template.ErrTemplateNameEmpty.Error())
}

func TestNoVaribleTemplate(t *testing.T) {
	tmpl := template.NewTemplateEngine("./test_data")
	resulte, err := tmpl.Render("novar", nil)

	assert.NoError(t, err, "expected no error")
	assert.Equal(t, "test\n", resulte, "expect test got: ")
}

func TestTemplateNoteFound(t *testing.T) {
	tmpl := template.NewTemplateEngine("./test_data")
	_, err := tmpl.Render("notfound", nil)

	assert.ErrorContains(t, err, template.ErrTemplateNotFound.Error())
}

func TestTemplateWithVar(t *testing.T) {
	tmpl := template.NewTemplateEngine("./test_data")
	resulte, err := tmpl.Render("withvar", map[string]string{
		"key": "testkey",
	})

	assert.NoError(t, err)
	assert.Equal(t, "key: testkey\n", resulte)
}
