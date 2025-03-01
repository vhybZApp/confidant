package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigFunc func(*Config)

type Config struct {
	AzurOpenAIConf *AzurOpenAIConfig
	TemplatePath   string
}

type AzurOpenAIConfig struct {
	Key string
	URL string
}

func defaultConfig() Config {
	return Config{
		AzurOpenAIConf: &AzurOpenAIConfig{
			Key: "",
			URL: "",
		},
		TemplatePath: "./tmpl",
	}
}

func Configuration(configs ...ConfigFunc) *Config {
	config := defaultConfig()
	for _, f := range configs {
		f(&config)
	}
	return &config
}

func WithDotEnvConfig(conf *Config) {
	envFile := ".env"
	if IsTestEnv() {
		envFile = "../../.env.test"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("loading .env file: %v", err)
	}

	key := os.Getenv("AZURE_OPENAI_KEY")
	if key != "" {
		conf.AzurOpenAIConf.Key = key
	}
	url := os.Getenv("AZURE_OPENAI_URL")
	if url != "" {
		conf.AzurOpenAIConf.URL = url
	}

	tmpl := os.Getenv("TEMPLATE_PATH")
	if tmpl != "" {
		conf.TemplatePath = tmpl
	}
}

func IsTestEnv() bool {
	return flag.Lookup("test.v") != nil
}
