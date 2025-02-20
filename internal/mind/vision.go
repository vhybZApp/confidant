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

type BoundedBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}
type Annotation struct {
	BoundedBox BoundedBox `json:"bbox"`
	CType      string     `json:"type"`
	Contetnt   string     `jsnon:"content"`
	ID         int        `json:"id"`
}

type AnnotatedImage struct {
	ImageBase64 string       `json:"image_base64"`
	Annotations []Annotation `json:"annotations"`
	ScreenInfo  string       `json:"screen_info"`
}

type Vision struct {
	client *omni.Client
}

func (v Vision) Annotate(description string, input io.Reader) (AnnotatedImage, error) {
	ai := AnnotatedImage{}

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
			BoundedBox: BoundedBox{
				X:      item.BBox[0],
				Y:      item.BBox[1],
				Height: item.BBox[2],
				Width:  item.BBox[3],
			},
			CType:    item.Type,
			Contetnt: item.Content,
			ID:       i,
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
