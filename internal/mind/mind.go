package mind

import (
	"fmt"
	"log"

	"github.com/farhoud/confidant/internal/template"
	"github.com/farhoud/confidant/pkg/omni"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type mindService struct {
	ready  bool
	llm    *LLM
	vision *Vision
	tmpl   template.Template
	screen Inspect
}

func (m mindService) Ready() bool {
	return m.ready
}

func (m mindService) Plan(goal string) ([]Action, error) {
	plan := []Action{}

	sm, err := m.tmpl.Render("planner-system", nil)
	if err != nil {
		return plan, fmt.Errorf("unable to render template: %w", err)
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(sm),
	}

	for {
		reader, err := m.screen.Inspect()
		if err != nil {
			return plan, err
		}

		andi, err := m.vision.Annotate("screen", reader)
		if err != nil {
			return plan, err
		}

		tmv := map[string]string{
			"ScreenInfo": andi.ScreenInfo,
		}

		if len(messages) == 1 {
			tmv["Goal"] = goal
		}

		um, err := m.tmpl.Render("planner-user", tmv)
		if err != nil {
			return plan, fmt.Errorf("unable to render template: %w", err)
		}

		msg_content := []openai.ChatCompletionContentPartUnionParam{
			openai.ChatCompletionContentPartTextParam{Text: openai.F(um), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)},
		}

		dataURl := DataURL("image/png", andi.ImageBase64)
		messages = append(messages, openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F(append(msg_content, openai.ImagePart(dataURl))),
		})

		msg, err := m.llm.Call(messages)
		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return plan, err
		}
		messages = append(messages, msg)
		action, err := ParseLLMActionResponse(msg.Content)
		if err != nil {
			return plan, err
		}

		plan = append(plan, action)
		if action.NextAction == "None" {
			break
		}

	}
	fmt.Printf("%#v", plan)
	return plan, nil
}

func NewMind(url, token string, screen Inspect) *mindService {
	if url == "" || token == "" {
		return &mindService{ready: false}
	}
	oc := openai.NewClient(
		option.WithBaseURL(url),
		option.WithAPIKey(token),
	)
	omni := omni.NewClient("http://localhost:8000")
	tmpl := template.NewTemplateEngine("./tmpl")

	llm := NewLLM(oc, tmpl, "azure-gpt-4o")
	vision := NewVision(omni)

	return &mindService{
		ready:  true,
		llm:    &llm,
		tmpl:   tmpl,
		vision: &vision,
		screen: screen,
	}
}

type Action struct {
	Reasoning  string `json:"Reasoning"`
	NextAction string `json:"Next Action"`
	BoxID      int    `json:"Box ID"`
	Value      string `json:"value"`
}
