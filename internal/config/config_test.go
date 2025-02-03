package config_test

import (
	"testing"

	"github.com/farhoud/confidant/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	conf := config.Configuration()

	assert.Equal(t, conf.AzurOpenAIConf.Key, "")
	assert.Equal(t, conf.AzurOpenAIConf.URL, "")
}

func TestDotEnvFile(t *testing.T) {
	conf := config.Configuration(config.WithDotEnvConfig)

	assert.Equal(t, conf.AzurOpenAIConf.Key, "testkey")
	assert.Equal(t, conf.AzurOpenAIConf.URL, "testurl")
	assert.Equal(t, conf.TemplatePath, "./tmpl")
}
