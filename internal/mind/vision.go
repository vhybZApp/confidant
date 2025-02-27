package mind

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/farhoud/confidant/pkg/omni"
)

var ErrBlindVision = errors.New("input is null we are blind")

type Annotation struct {
	BoundedBox [4]float64 `json:"bbox"`
	CType      string     `json:"type"`
	Contetnt   string     `jsnon:"content"`
	ID         int        `json:"id"`
}

type AnnotatedImage struct {
	ImageBase64 string       `json:"image_base64"`
	Annotations []Annotation `json:"annotations"`
	ScreenInfo  string       `json:"screen_info"`
	Width       int
	Height      int
}

func (a AnnotatedImage) BoundedBox(id int) (int, int) {
	annotation := a.Annotations[id]

	x := ((annotation.BoundedBox[2] - annotation.BoundedBox[0]) / 2) + annotation.BoundedBox[0]
	y := ((annotation.BoundedBox[3] - annotation.BoundedBox[1]) / 2) + annotation.BoundedBox[1]
	return int(x * float64(a.Width)), int(y * float64(a.Height))
}

type Vision struct {
	client *omni.Client
}

func (v Vision) Annotate(description string, screenSize []int, input io.Reader) (AnnotatedImage, error) {
	ai := AnnotatedImage{}

	ai.Width = screenSize[0]
	ai.Height = screenSize[1]

	ib64, err := EncodeToBase64(input)
	if err != nil {
		return ai, err
	}

	or, err := v.client.Parse(context.TODO(), ib64)
	if err != nil {
		return ai, err
	}

	ai.Annotations = make([]Annotation, len(or.ParsedContentList))
	ai.ImageBase64 = or.ImageBase64
	for i, item := range or.ParsedContentList {
		switch item.Type {
		case "text":
			ai.ScreenInfo += fmt.Sprintf("ID: %d, Text: %s \n", i, item.Content)
		case "icon":
			ai.ScreenInfo += fmt.Sprintf("ID: %d, Icon: %s \n", i, item.Content)
		}
		ai.Annotations[i] = Annotation{
			BoundedBox: item.BBox,
			CType:      item.Type,
			Contetnt:   item.Content,
			ID:         i,
		}
	}
	return ai, nil
}

func NewVision(client *omni.Client) Vision {
	return Vision{client: client}
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
