package template_test

import (
	"testing"

	"github.com/farhoud/confidant/internal/template"
	"github.com/stretchr/testify/assert"
)

func TestEmptyNameRender(t *testing.T) {
	_, err := template.Render("", nil)

	assert.EqualError(t, err, template.ErrTemplateNameEmpty.Error())
}

func TestNoVaribleTemplate(t *testing.T) {
	resulte, err := template.Render("novar", nil)

	assert.NoError(t, err, "expected no error")
	assert.Equal(t, "test\n", resulte, "expect test got: ")
}
