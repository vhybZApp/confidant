package mind

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAction(t *testing.T) {
	content := "<think>\nThe task is to open Safari. Looking at the screenshot, I can see various icons in the dock at the bottom. Identifying the correct icon for Safari is crucial. The icon in ID 14 in the dock is marked as \"Safari.\"\n\nTo achieve the task, I need to click on the Safari icon in the dock.\n</think>\n<output>\n{\n    \"Reasoning\": \"To open Safari, I need to click on the Safari icon in the dock. The icon for Safari is marked with ID 14.\",\n    \"Next Action\": \"left_click\",\n  \"Box ID\": 14\n}\n</output>"

	action, err := ParseLLMActionResponse(content)
	t.Logf("plan: %v", action)
	t.Logf("err: %v", err)
	assert.NoError(t, err, "no error should happend")
	assert.NotEmpty(t, action, "plan should not be empty")
}
