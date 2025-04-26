package mind

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/farhoud/confidant/internal/config"
	"github.com/farhoud/confidant/internal/template"
	"github.com/farhoud/confidant/pkg/omni"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type runnerService struct {
	ready      bool
	conf       config.Config
	planner    Agent
	operator   Agent
	controller Agent
}

func (m runnerService) Ready() bool {
	return m.ready
}

func (m runnerService) Run(goal string) error {
	if !m.ready {
		return errors.New("mind is not ready")
	}
	thread := NewThread(rand.Intn(10000))
	m.planner.Achieve(goal, thread)
	for {
		m.controller.Achieve(goal, thread)
		for {
			if thread.Acheived {
				break
			}
			m.operator.Achieve(goal, thread)
			time.Sleep(time.Second * 5)
		}
		thread.Acheived = false
	}
	return nil
}

func NewApp(conf config.Config, screen Inspect) *runnerService {
	if conf.AzurOpenAIConf.URL == "" || conf.AzurOpenAIConf.Key == "" {
		return &runnerService{ready: false}
	}
	oc := openai.NewClient(
		option.WithBaseURL(conf.AzurOpenAIConf.URL),
		option.WithAPIKey(conf.AzurOpenAIConf.Key),
	)
	omni := omni.NewClient("http://localhost:8000")
	tmpl := template.NewTemplateEngine(conf.TemplatePath)

	llm := NewLLM(oc, conf.LLMModel)
	vision := NewVision(omni)
	p := NewPlanner(&llm, tmpl, screen, conf.DeviceType)
	c := NewController(&llm, tmpl, screen, conf.DeviceType)
	o := NewOperator(&llm, tmpl, screen, vision, conf.DeviceType)
	return &runnerService{
		ready:      true,
		conf:       conf,
		planner:    p,
		operator:   o,
		controller: c,
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
