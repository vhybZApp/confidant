package mind

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/farhoud/confidant/internal/template"
	"github.com/go-vgo/robotgo"
	"github.com/openai/openai-go"
)

type operator struct {
	tmpl   template.Template
	screen Inspect
	vision Vision
	divt   string
	llm    *LLM
}

func (o operator) Achieve(goal string, thread *Thread) error {
	reader, err := o.screen.Inspect()
	if err != nil {
		return err
	}
	sw, sh := robotgo.GetScreenSize()
	andi, err := o.vision.Annotate("screen", []int{sw, sh}, reader)
	if err != nil {
		return err
	}

	tmv := map[string]interface{}{
		"ScreenInfo":  andi.ScreenInfo,
		"Goal":        goal,
		"ImageBase64": andi.ImageBase64,
		"DeviceType":  o.divt,
	}

	mem, err := o.Match(goal, tmv, thread)
	msg, err := o.llm.Call(mem)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}

	mem = append(mem, openai.AssistantMessage(msg.Content))
	fmt.Printf("Assistant message: %s", msg.Content)
	action, err := ParseLLMActionResponse(msg.Content)
	if err != nil {
		return err
	}

	s := Snapshot{
		Agent:    "operator",
		Messages: mem,
	}
	reader.Seek(0, io.SeekStart)
	s.AddAttachment(reader)
	s.AddAttachmentFromBase64(andi.ImageBase64)
	thread.AddSnapshot(s)
	thread.Store("./data")
	if action.NextAction == "None" {
		thread.GoalAcheived()
		return nil
	}

	err = ExecAction(action, andi)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	return nil
}

func (o operator) Match(goal string, inputInfo map[string]interface{}, thread *Thread) ([]openai.ChatCompletionMessageParamUnion, error) {
	sm, err := o.tmpl.Render("operator-system", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	um, err := o.tmpl.Render("operator-user", inputInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	mem := []openai.ChatCompletionMessageParamUnion{}
	s := thread.LatestSnapShot("operator").Messages
	// plan := thread.LatestSnapShot("planner").Messages[len(thread.LatestSnapShot("planner").Messages)-1]
	if len(s) < 1 {
		mem = append(mem, openai.SystemMessage(sm))
	} else {
		for i, item := range s {
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

func NewOperator(llm *LLM, tmpl template.Template, screen Inspect, vision Vision, deviceType string) *planner {
	return &planner{
		llm:    llm,
		tmpl:   tmpl,
		screen: screen,
		divt:   deviceType,
	}
}
