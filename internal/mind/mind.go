package mind

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/farhoud/confidant/internal/template"
	"github.com/farhoud/confidant/pkg/omni"

	"github.com/openai/openai-go"
)

type mindService struct {
	ready  bool
	client *openai.Client
	tmpl   template.Template
	op     *omni.Client
}

func (m mindService) Ready() bool {
	return m.ready
}

func (m mindService) Detect(d string, i io.ReadSeeker) (Box, error) {
	if i == nil {
		return Box{}, ErrBlindVision
	}

	client := m.client

	ib64, err := EncodeToBase64(i)
	if err != nil {
		return Box{}, err
	}

	or, err := m.op.Parse(context.TODO(), ib64)
	if err != nil {
		return Box{}, err
	}

	dataURl := DataURL("image/png", or.ImageBase64)
	sm, err := m.tmpl.Render("vision-detection-system", or)
	if err != nil {
		return Box{}, err
	}

	log.Printf("sm: %s", sm)

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sm),
			openai.ChatCompletionUserMessageParam{
				Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
				Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ImagePart(dataURl)}),
			},
		}),
		Model: openai.F("azure-gpt-4o"),
	})
	if err != nil {
		return Box{}, nil
	}

	box := Box{}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &box)
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		log.Printf("resp: %v", resp.Choices[0].Message.Content)
		return box, err
	}
	return box, nil
}

func (m mindService) Plan() (Plan, error) {
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
