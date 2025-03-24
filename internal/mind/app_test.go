package mind

//
// import (
// 	"io"
// 	"testing"
//
// 	"github.com/farhoud/confidant/internal/config"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestOpenAIMindWithEmptyCredential(t *testing.T) {
// 	url := ""
// 	token := ""
// 	mind := NewMind(url, token, "", "", "", nil)
// 	assert.False(t, mind.Ready(), "when new mind create with empty url and token it should not ready")
// }
//
// func TestOpenMindAIWithValidCredentail(t *testing.T) {
// 	url := "https://api.openai.com/v1/completions"
// 	token := "sk-"
// 	mind := NewMind(url, token, "", "", "", nil)
// 	assert.True(t, mind.Ready(), "when new mind create with valid url and token it should ready")
// }
//
// func TestOpenAIMindWithEmptyRequest(t *testing.T) {
// 	conf := config.Configuration(config.WithDotEnvConfig)
// 	assert.NotEmpty(t, conf.AzurOpenAIConf.Key)
// 	assert.NotEmpty(t, conf.AzurOpenAIConf.URL)
// 	mind := NewMind(conf.AzurOpenAIConf.URL, conf.AzurOpenAIConf.Key, conf.TemplatePath, conf.LLMModel, conf.DeviceType, nil)
// 	assert.True(t, mind.Ready(), "new mind should be ready")
// }
//
// func TestMVI(t *testing.T) {
// 	screens := []string{"test_data/mac-desktop.jpg", "test_data/mac_safari.png"}
// 	mvi := NewMockScreenInspector(screens)
// 	for i := 0; i < len(screens); i++ {
// 		r, err := mvi.Inspect()
// 		assert.NoError(t, err, "inspect should not return error")
// 		data, err := io.ReadAll(r)
// 		assert.NoError(t, err, "no error should happen")
// 		assert.NotEmpty(t, data, "not empty")
// 	}
// }
//
// func TestMockScreenOpenSafari(t *testing.T) {
// 	conf := config.Configuration(config.WithDotEnvConfig)
//
// 	mvi := NewMockScreenInspector([]string{"test_data/mac-desktop.jpg", "test_data/mac_safari.png"})
// 	mind := NewMind(conf.AzurOpenAIConf.URL, conf.AzurOpenAIConf.Key, conf.TemplatePath, conf.LLMModel, conf.DeviceType, mvi)
//
// 	plan, err := mind.Plan("open safari")
//
// 	t.Logf("plan: %v", plan)
// 	t.Logf("err: %v", err)
// 	assert.NoError(t, err, "no error should happend")
// 	assert.NotEmpty(t, plan, "plan should not be empty")
// }
//
// func TestOpenChrome(t *testing.T) {
// 	conf := config.Configuration(config.WithDotEnvConfig)
//
// 	mvi := NewRobotScreenInspector()
// 	mind := NewMind(conf.AzurOpenAIConf.URL, conf.AzurOpenAIConf.Key, conf.TemplatePath, conf.LLMModel, conf.DeviceType, mvi)
//
// 	plan, err := mind.Plan("open safari")
//
// 	t.Logf("plan: %v", plan)
// 	t.Logf("err: %v", err)
// 	assert.NoError(t, err, "no error should happend")
// 	assert.NotEmpty(t, plan, "plan should not be empty")
// }
