package mind

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

// EncodeToBase64 reads from an io.Reader, detects the MIME type, and encodes the content to a Base64 image string.
func EncodeToBase64(reader io.Reader) (string, error) {
	// Read the file into a buffer
	buf := &bytes.Buffer{}
	tee := io.TeeReader(reader, buf)

	// Read a small portion to detect the MIME type
	sniff := make([]byte, 512) // First 512 bytes are enough for MIME detection
	n, err := tee.Read(sniff)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Detect MIME type
	mimeType := http.DetectContentType(sniff[:n])

	// Ensure it's a valid image type
	if !strings.HasPrefix(mimeType, "image/") {
		return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	// Read full content from buffer
	data, err := io.ReadAll(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Encode to Base64
	base64Str := base64.StdEncoding.EncodeToString(data)

	// Construct Data URL format
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)

	log.Printf("dataURL: %s", dataURL)
	return dataURL, nil
}
