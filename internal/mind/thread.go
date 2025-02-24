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

type Snapshot = []openai.ChatCompletionMessageParamUnion

type Thread struct {
	ID          int         `json:"id"`
	Attachments []io.Reader `json:"Attachments"`
	History     []Snapshot  `json:"History"`
	mutex       sync.Mutex  // To handle concurrent writes
}

// NewThread creates a new thread instance
func NewThread(id int) *Thread {
	return &Thread{
		ID:          id,
		Attachments: []io.Reader{},
		History:     []Snapshot{},
	}
}

func (t *Thread) AddAttachmentFromBase64(encoded string) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %w", err)
	}
	buffer := bytes.NewBuffer(data)
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Attachments = append(t.Attachments, buffer)
	return nil
}

// AddAttachment adds an attachment to the thread
func (t *Thread) AddAttachment(attachment io.Reader) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
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

	for i, attachment := range t.Attachments {
		filePath := filepath.Join(attachmentsDir, fmt.Sprintf("attachment_%d", i))
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
	}

	return nil
}
