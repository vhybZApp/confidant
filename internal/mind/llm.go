package mind

import (
	"context"
	"encoding/json"
	"log"

	"github.com/farhoud/confidant/internal/template"
	"github.com/openai/openai-go"
)

type LLM struct {
	client *openai.Client
	model  string
}

func (l LLM) Call(messages []openai.ChatCompletionMessageParamUnion, v any) error {
	resp, err := l.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(l.model),
	})
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), v)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		log.Printf("resp: %v", resp.Choices[0].Message.Content)
		return err
	}

	messages = append(messages, resp.Choices[0].Message)

	return nil
}

func NewLLM(client *openai.Client, tmpl template.Template, model string) LLM {
	return LLM{
		client: client,
		model:  model,
	}
}
