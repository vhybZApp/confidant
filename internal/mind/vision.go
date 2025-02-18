package mind

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	Detect(description string, input io.ReadSeeker) (Box, error)
}

func EncodeToBase64(reader io.Reader) (string, error) {
	// Read full content from buffer
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Encode to Base64
	base64Str := base64.StdEncoding.EncodeToString(data)

	// Construct Data URL format
	// dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
	dataURL := base64Str

	return dataURL, nil
}

// Public function to get MIME type from an io.Reader
func MimeType(r io.Reader) (string, error) {
	// Read the first 512 bytes to determine MIME type
	buffer := make([]byte, 512)
	_, err := r.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Get the MIME type using the buffer
	return http.DetectContentType(buffer), nil
}

// Public function to generate a Data URL
func DataURL(mimeType, encodedData string) string {
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encodedData)
}
