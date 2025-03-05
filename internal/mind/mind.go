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
	screen Inspect
	mem    *Memory
	divt   string
}

func (m mindService) Ready() bool {
	return m.ready
}

func (m mindService) Plan(goal string) ([]Action, error) {
	thread := NewThread(rand.Intn(10000))
	plan := []Action{}

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

		tmv := map[string]interface{}{
			"ScreenInfo":  andi.ScreenInfo,
			"Goal":        goal,
			"ImageBase64": andi.ImageBase64,
			"DeviceType":  m.divt,
		}

		mem, err := m.mem.Match(goal, tmv, thread.History)
		msg, err := m.llm.Call(mem)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return plan, err
		}

		mem = append(mem, openai.AssistantMessage(msg.Content))
		fmt.Printf("Assistant message: %s", msg.Content)
		action, err := ParseLLMActionResponse(msg.Content)
		if err != nil {
			return plan, err
		}

		thread.AddSnapshot(mem)
		reader.Seek(0, io.SeekStart)
		thread.AddAttachment(reader)
		thread.AddAttachmentFromBase64(andi.ImageBase64)
		thread.Store("./data")
		plan = append(plan, action)
		if action.NextAction == "None" {
			break
		}

		err = ExecAction(action, andi)
		if err != nil {
			return plan, err
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Printf("%#v", plan)
	return plan, nil
}

func NewMind(url, token, tmplPath, llmModel, deviceType string, screen Inspect) *mindService {
	if url == "" || token == "" {
		return &mindService{ready: false}
	}
	oc := openai.NewClient(
		option.WithBaseURL(url),
		option.WithAPIKey(token),
	)
	omni := omni.NewClient("http://localhost:8000")
	tmpl := template.NewTemplateEngine(tmplPath)
	mem := NewMemory(tmpl)

	llm := NewLLM(oc, llmModel)
	vision := NewVision(omni)

	return &mindService{
		ready:  true,
		llm:    &llm,
		mem:    &mem,
		vision: &vision,
		screen: screen,
		divt:   deviceType,
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
