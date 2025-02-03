package mind

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type mind struct {
	ready  bool
	client *azopenai.Client
}

func (m mind) Ready() bool {
	return m.ready
}

func (m mind) Plan() (Plan, error) {
	return Plan{
		Actions: []Action{},
	}, nil
}

type Plan struct {
	Actions []Action
}

type Action struct {
	Expect string
	Output string
	Func   string
	Args   []string
}

func NewMind(url, key string) *mind {
	m := &mind{}

	if url == "" || key == "" {
		return m
	}

	keyCredential := azcore.NewKeyCredential(key)

	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
	client, err := azopenai.NewClientWithKeyCredential(url, keyCredential, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		m.ready = false
		return m
	}

	m.ready = true
	m.client = client

	return m
}
