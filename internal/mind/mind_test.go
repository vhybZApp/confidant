package mind

import (
	"testing"

	"github.com/farhoud/confidant/internal/config"
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

func TestOpenAIMindWithEmptyRequest(t *testing.T) {
	conf := config.Configuration(config.WithDotEnvConfig)
	mind := NewMind(conf.AzurOpenAIConf.URL, conf.AzurOpenAIConf.Key)
	assert.True(t, mind.Ready(), "new mind should be ready")
}

func TestSimplePlanning(t *testing.T) {
	conf := config.Configuration(config.WithDotEnvConfig)
	mind := NewMind(conf.AzurOpenAIConf.URL, conf.AzurOpenAIConf.Key)

	plan, err := mind.Plan()

	assert.NoError(t, err, "no error should happend")
	assert.NotEmpty(t, plan, "plan should not be empty")
}
