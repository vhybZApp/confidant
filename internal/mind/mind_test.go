package mind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAIMindWithEmptyCredential(t *testing.T) {
	url := ""
	token := ""
	mind := NewMind(url, token)
	assert.False(t, mind.Ready(), "when new mind create with empty url and token it should not ready")
}

func TestOpenMindAIWithValidCredentail(t *testing.T) {
	url := "https://api.openai.com/v1/completions"
	token := "sk-"
	mind := NewMind(url, token)
	assert.True(t, mind.Ready(), "when new mind create with valid url and token it should ready")
}
