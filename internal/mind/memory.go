package mind

import (
	"fmt"

	"github.com/farhoud/confidant/internal/template"
	"github.com/openai/openai-go"
)

type Memory struct {
	tmpl template.Template
}

func (m Memory) Match(goal string, inputInfo map[string]interface{}, s []Snapshot) ([]openai.ChatCompletionMessageParamUnion, error) {
	sm, err := m.tmpl.Render("planner-system", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	um, err := m.tmpl.Render("planner-user", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	mem := []openai.ChatCompletionMessageParamUnion{}
	if len(s) < 1 {
		mem = append(mem, openai.SystemMessage(sm))
	} else {
		for i, item := range s[len(s)-1] {
			switch item.(type) {
			case openai.ChatCompletionSystemMessageParam:
				mem = append(mem, openai.SystemMessage(sm))
			case openai.ChatCompletionUserMessageParam:
				if i == 1 {
					continue
				}
				mem = append(mem, openai.UserMessage("action done"))
			default:
				mem = append(mem, item)
			}
		}
	}
	msg_content := []openai.ChatCompletionContentPartUnionParam{
		openai.ChatCompletionContentPartTextParam{Text: openai.F(um), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)},
	}

	dataURl := DataURL("image/png", inputInfo["ImageBase64"].(string))
	mem = append(mem, openai.ChatCompletionUserMessageParam{
		Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
		Content: openai.F(append(msg_content, openai.ImagePart(dataURl))),
	})
	return mem, nil
}

func NewMemory(tmpl template.Template) Memory {
	return Memory{tmpl: tmpl}
}
