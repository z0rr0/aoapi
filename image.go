package aoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ImageSize is a type of image size.
type ImageSize string

// Image sizes.
const (
	ImageSize256  ImageSize = "256x256"
	ImageSize512  ImageSize = "512x512"
	ImageSize1024 ImageSize = "1024x1024"
)

// ImageRequest is a struct of image request.
type ImageRequest struct {
	Prompt string    `json:"prompt"`
	N      uint      `json:"n,omitempty"`
	Size   ImageSize `json:"size,omitempty"`
}

func (i *ImageRequest) marshal() (io.Reader, error) {
	if i.Prompt == "" {
		return nil, errors.Join(ErrRequiredParam, fmt.Errorf("prompt must not be empty"))
	}

	data, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal image request: %w", err)
	}

	return bytes.NewReader(data), nil
}

func (i *ImageRequest) build(ctx context.Context, auth *Params) (*http.Request, error) {
	body, err := i.marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, auth.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create image request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.Bearer))

	if auth.Organization != "" {
		req.Header.Set("OpenAI-Organization", auth.Organization)
	}

	return req, nil
}

// ImageData stores image URL.
type ImageData struct {
	URL string `json:"url"`
}

// ImageResponse is a struct of image response.
type ImageResponse struct {
	Created   int64       `json:"created"`
	Data      []ImageData `json:"data"`
	CreatedTs time.Time   `json:"-"`
}

func (ir *ImageResponse) build(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(&ir); err != nil {
		return fmt.Errorf("failed to unmarshal image response: %w", err)
	}

	if len(ir.Data) == 0 {
		return errors.Join(ErrResponse, fmt.Errorf("empty image response"))
	}

	ir.CreatedTs = time.Unix(ir.Created, 0)
	return nil
}

// String returns a string representation of the image response.
func (ir *ImageResponse) String() string {
	const sep = "\n"
	var buf bytes.Buffer

	for i, d := range ir.Data {
		buf.WriteString(fmt.Sprintf("%d. ", i+1))
		buf.WriteString(d.URL)
		buf.WriteString(sep)
	}

	return buf.String()
}

// Image sends request to the image API.
func Image(ctx context.Context, client *http.Client, i *ImageRequest, p Params) (*ImageResponse, error) {
	body, err := commonRequest(ctx, client, i, p)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = body.Close()
	}()

	response := &ImageResponse{}
	if err = response.build(body); err != nil {
		return nil, err
	}

	return response, nil
}
