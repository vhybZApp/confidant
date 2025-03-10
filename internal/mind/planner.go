package mind

import (
	"fmt"
	"io"
	"log"

	"github.com/farhoud/confidant/internal/template"
	"github.com/openai/openai-go"
)

type planner struct {
	ready  bool
	llm    *LLM
	screen Inspect
	divt   string
	tmpl   template.Template
}

func (c planner) Achieve(goal string, thread *Thread) error {
	reader, err := c.screen.Inspect()
	if err != nil {
		return err
	}

	bs64i, err := EncodeToBase64(reader)
	if err != nil {
		return err
	}

	tmv := map[string]interface{}{
		"Goal":        goal,
		"ImageBase64": bs64i,
		"DeviceType":  c.divt,
	}

	mem, err := c.Match(goal, tmv, thread.History)
	msg, err := c.llm.Call(mem)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}

	mem = append(mem, openai.AssistantMessage(msg.Content))
	fmt.Printf("Assistant message: %s", msg.Content)

	s := Snapshot{
		Agent:    "planner",
		Messages: mem,
	}
	reader.Seek(0, io.SeekStart)
	s.AddAttachment(reader)
	thread.AddSnapshot(s)
	thread.Store("./data")
	return nil
}

func (c planner) Match(goal string, inputInfo map[string]interface{}, s []Snapshot) ([]openai.ChatCompletionMessageParamUnion, error) {
	sm, err := c.tmpl.Render("planner-system", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	um, err := c.tmpl.Render("planner-user", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	mem := []openai.ChatCompletionMessageParamUnion{}
	mem = append(mem, openai.SystemMessage(sm))
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

func NewPlanner(llm *LLM, tmpl template.Template, screen Inspect, deviceType string) *planner {
	return &planner{
		llm:    llm,
		tmpl:   tmpl,
		screen: screen,
		divt:   deviceType,
	}
}
