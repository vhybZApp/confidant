package mind

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/farhoud/confidant/internal/template"
	"github.com/openai/openai-go"
)

type LLM struct {
	client *openai.Client
	model  string
}

func (l LLM) Call(messages []openai.ChatCompletionMessageParamUnion) (openai.ChatCompletionMessage, error) {
	resp, err := l.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(l.model),
	})
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}

	// fmt.Printf("llm context: %v \n llm response %v", messages, resp.Choices[0].Message)

	return resp.Choices[0].Message, nil
}

func NewLLM(client *openai.Client, tmpl template.Template, model string) LLM {
	return LLM{
		client: client,
		model:  model,
	}
}

func ParseLLMActionResponse(text string) (Action, error) {
	action := Action{}
	fmt.Printf("llm response: %v", text)
	re := regexp.MustCompile(`(?s)<output>\s*(\{.*?\})\s*</output>`)
	matches := re.FindStringSubmatch(text)

	if len(matches) > 1 {
		jsonStr := matches[1]

		// Struct to hold extracted JSON data
		err := json.Unmarshal([]byte(jsonStr), &action)
		if err != nil {
			return action, err
		}
	} else {
		return action, errors.New("extract failed")
	}

	return action, nil
}
