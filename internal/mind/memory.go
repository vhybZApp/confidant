package mind

import (
	"github.com/openai/openai-go"
)

func Revise(goal string, s []openai.ChatCompletionMessageParamUnion) []openai.ChatCompletionMessageParamUnion {
	mem := []openai.ChatCompletionMessageParamUnion{}
	for i, item := range s {
		switch item.(type) {
		case openai.ChatCompletionUserMessageParam:
			if i == 1 {
				mem = append(mem, openai.UserMessage("goal: "+goal))
			}
			mem = append(mem, openai.UserMessage("action done"))
		default:
			mem = append(mem, item)
		}
	}
	return mem
}
