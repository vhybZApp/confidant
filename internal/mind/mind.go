package mind

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/farhoud/confidant/internal/template"
	"github.com/farhoud/confidant/pkg/omni"
	"github.com/go-vgo/robotgo"
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
	thread := NewThread(rand.Intn(10000))
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
		sw, sh := robotgo.GetScreenSize()
		andi, err := m.vision.Annotate("screen", []int{sw, sh}, reader)
		if err != nil {
			return plan, err
		}

		tmv := map[string]string{
			"ScreenInfo": andi.ScreenInfo,
			"Goal":       goal,
		}

		if len(messages) == 1 {
			tmv["Goal"] = goal
		} else {
			messages = Revise(goal, messages)
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
			log.Printf("ERROR: %s", err)
			return plan, err
		}
		messages = append(messages, openai.AssistantMessage(msg.Content))
		action, err := ParseLLMActionResponse(msg.Content)
		if err != nil {
			return plan, err
		}

		thread.AddSnapshot(messages)
		reader.Seek(0, io.SeekStart)
		thread.AddAttachment(reader)
		thread.AddAttachmentFromBase64(andi.ImageBase64)
		thread.Store("./data")
		plan = append(plan, action)
		fmt.Printf("action: %+v", action)
		if action.NextAction == "None" {
			break
		}

		err = ExecAction(action, andi)
		if err != nil {
			return plan, err
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Printf("%#v", plan)
	return plan, nil
}

func NewMind(url, token, tmplPath string, screen Inspect) *mindService {
	if url == "" || token == "" {
		return &mindService{ready: false}
	}
	oc := openai.NewClient(
		option.WithBaseURL(url),
		option.WithAPIKey(token),
	)
	omni := omni.NewClient("http://localhost:8000")
	tmpl := template.NewTemplateEngine(tmplPath)

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

func (a Action) IntValue() (int, error) {
	if a.Value == "" {
		return 0, errors.New("no value")
	}
	return strconv.Atoi(a.Value)
}
