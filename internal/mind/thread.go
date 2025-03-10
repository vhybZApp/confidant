package mind

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/openai/openai-go"
)

type Snapshot struct {
	Attachments []io.Reader `json:"Attachments"`
	Messages    []openai.ChatCompletionMessageParamUnion
	Agent       string
}

type Thread struct {
	ID       int        `json:"id"`
	History  []Snapshot `json:"History"`
	Acheived bool
	mutex    sync.Mutex // To handle concurrent writes
}

// NewThread creates a new thread instance
func NewThread(id int) *Thread {
	return &Thread{
		ID:      id,
		History: []Snapshot{},
	}
}

func (t *Thread) GoalAcheived() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Acheived = true
}

func (t *Thread) LatestSnapShot(agent string) *Snapshot {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	i := len(t.History) - 1
	for {
		if t.History[i].Agent == agent {
			return &t.History[i]
		}
		i--
		if i < 0 {
			return nil
		}
	}
}

func (t *Snapshot) AddAttachmentFromBase64(encoded string) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %w", err)
	}
	buffer := bytes.NewBuffer(data)
	t.Attachments = append(t.Attachments, buffer)
	return nil
}

// AddAttachment adds an attachment to the thread
func (t *Snapshot) AddAttachment(attachment io.Reader) {
	t.Attachments = append(t.Attachments, attachment)
}

// AddSnapshot adds a snapshot to the thread history
func (t *Thread) AddSnapshot(snapshot Snapshot) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.History = append(t.History, snapshot)
}

// Store saves the thread data in a directory named after the thread ID.
func (t *Thread) Store(basePath string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	dirPath := filepath.Join(basePath, fmt.Sprintf("%d", t.ID))
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Store snapshots
	snapshotFile := filepath.Join(dirPath, "history.json")
	existingSnapshots := make([]Snapshot, 0)

	if data, err := os.ReadFile(snapshotFile); err == nil {
		_ = json.Unmarshal(data, &existingSnapshots)
	}

	if len(t.History) > len(existingSnapshots) {
		// Append only new snapshots
		newSnapshots := t.History[len(existingSnapshots):]
		allSnapshots := append(existingSnapshots, newSnapshots...)
		data, err := json.MarshalIndent(allSnapshots, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal snapshots: %w", err)
		}
		if err := os.WriteFile(snapshotFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write history file: %w", err)
		}
	}

	// Store attachments
	attachmentsDir := filepath.Join(dirPath, "attachments")
	if err := os.MkdirAll(attachmentsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create attachments directory: %w", err)
	}

	i := 0
	for _, snapshot := range t.History {
		for _, attachment := range snapshot.Attachments {
			filePath := filepath.Join(attachmentsDir, fmt.Sprintf("attachment_%s_%d", snapshot.Agent, i))
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				file, err := os.Create(filePath)
				if err != nil {
					return fmt.Errorf("failed to create attachment file: %w", err)
				}
				defer file.Close()
				if _, err := io.Copy(file, attachment); err != nil {
					return fmt.Errorf("failed to write attachment file: %w", err)
				}
			}
			i++
		}
	}

	return nil
}
