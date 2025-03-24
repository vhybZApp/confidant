package mind

import (
	"fmt"
	"io"
	"log"

	"github.com/farhoud/confidant/internal/template"
	"github.com/openai/openai-go"
)

type controller struct {
	tmpl   template.Template
	screen Inspect
	divt   string
	llm    *LLM
}

func (o controller) Achieve(goal string, thread *Thread) error {
	reader, err := o.screen.Inspect()
	if err != nil {
		return err
	}

	bs64, err := EncodeToBase64(reader)
	if err != nil {
		return err
	}
	tmv := map[string]interface{}{
		"Goal":        goal,
		"ImageBase64": bs64,
		"DeviceType":  o.divt,
	}

	mem, err := o.Match(goal, tmv, thread)
	msg, err := o.llm.Call(mem)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}

	mem = append(mem, openai.AssistantMessage(msg.Content))
	fmt.Printf("Controller message: %s", msg.Content)

	s := Snapshot{
		Agent:    "controller",
		Messages: mem,
	}
	reader.Seek(0, io.SeekStart)
	s.AddAttachment(reader)
	thread.AddSnapshot(s)
	thread.Store("./data")

	return nil
}

func (o controller) Match(goal string, inputInfo map[string]interface{}, thread *Thread) ([]openai.ChatCompletionMessageParamUnion, error) {
	sm, err := o.tmpl.Render("controller-system", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}

	sp := thread.LatestSnapShot("planner")
	msgs := sp.Messages
	if sp != nil && len(msgs) > 1 {
		inputInfo["Plan"] = msgs[len(msgs)-1].(openai.ChatCompletionAssistantMessageParam).Content
	}

	um, err := o.tmpl.Render("controller-user", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	mem := []openai.ChatCompletionMessageParamUnion{}

	sc := thread.LatestSnapShot("controller")
	so := thread.LatestSnapShot("operator")

	if sc == nil {
		mem = append(mem, openai.SystemMessage(sm))
		if sp != nil && len(sp.Messages) > 1 {
			mem = append(mem, openai.UserMessage(um))
		}
	} else {
		mem = append(mem, sc.Messages...)
		if so != nil && len(so.Messages) > 1 {
			mem = append(mem, so.Messages[1:]...)
		}
		mem = append(mem, openai.UserMessage("What is our next goal?"))
	}
	return mem, nil
}

func NewController(llm *LLM, tmpl template.Template, screen Inspect, deviceType string) *controller {
	return &controller{
		llm:    llm,
		tmpl:   tmpl,
		screen: screen,
		divt:   deviceType,
	}
}
