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

func (m mindService) Plan(goal string) (Plan, error) {
	plan := Plan{}
	reader, err := m.screen.Inspect()
	if err != nil {
		return plan, err
	}

	andi, err := m.vision.Annotate("screen", reader)
	if err != nil {
		return plan, err
	}
	sm, err := m.tmpl.Render("planner-system", nil)
	if err != nil {
		return plan, fmt.Errorf("unable to render template: %w", err)
	}

	log.Printf("sm: %s", sm)

	dataURl := DataURL("image/png", andi.ImageBase64)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(sm),
		openai.ChatCompletionUserMessageParam{
			Role: openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{
				openai.ImagePart(dataURl),
				openai.ChatCompletionContentPartTextParam{Text: openai.F(goal), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)},
			}),
		},
	}

	err = m.llm.Call(messages, &plan)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return plan, err
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

type Plan struct {
	Actions []Action `json:"Actions" jsonschema_description:"The actions needed to be done to achive the goal"`
}

type Action struct {
	Expect string   `json:"Expect" jsonschema_description:"What is the expected resulte of this action"`
	Output string   `json:"Output" jsonschema_description:"on which output should happend like mouse or keyboard"`
	Func   string   `json:"Func" jsonschema_description:"The function that should be called"`
	Args   []string `json:"Args" jsonschema_description:"The arguments that should be passed to the function"`
}
