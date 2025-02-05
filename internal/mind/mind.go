package mind

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/farhoud/confidant/internal/template"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Generate the JSON schema at initialization time
var HistoricalComputerResponseSchema = GenerateSchema[Plan]()

type mind struct {
	ready  bool
	client *openai.Client
	tmpl   template.Template
}

func (m mind) Ready() bool {
	return m.ready
}

func (m mind) Plan() (Plan, error) {
	client := m.client
	question := "Start browser"

	sm, err := m.tmpl.Render("planner-system", nil)
	if err != nil {
		return Plan{}, fmt.Errorf("unable to render template: %w", err)
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sm),
			openai.UserMessage(question),
		}),
		Model: openai.F("azure-gpt-4o"),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return Plan{}, nil
	}

	plan := Plan{}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &plan)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return plan, err
	}

	fmt.Printf("%#v", plan)
	return plan, nil
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

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func NewMind(url, key string) *mind {
	m := &mind{}

	if url == "" || key == "" {
		return m
	}

	m.tmpl = template.NewTemplateEngine("./tmpl")

	m.client = openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(url),
	)
	m.ready = true

	return m
}
