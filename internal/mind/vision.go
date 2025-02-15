package mind

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var ErrBlindVision = errors.New("input is null we are blind")

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Box struct {
	TRPoint Point `json:"top-left"`
	BLPoint Point `json:"bottom-right"`
}

type Vision interface {
	Detect(description string, input io.Reader) (Box, error)
}

// EncodeToBase64 reads an io.Reader and returns a Base64-encoded string
func EncodeToBase64(reader io.Reader) (string, error) {
	// Read the entire file into a buffer
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Encode to Base64
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str, nil
}

// GenerateOpenAIImageInput formats a Base64-encoded image for OpenAI's API
func GenerateOpenAIImageInput(base64Str, mimeType string) string {
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
}
