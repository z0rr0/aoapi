// Package aoapi provides a client for the OpenAI chat completion API.
package aoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	// ErrRequiredParam is an error that occurs when a required parameter is missing.
	ErrRequiredParam = errors.New("required parameter is missing")

	// ErrResponse is an error that occurs when the response is empty.
	ErrResponse = errors.New("failed response")

	// TokenLimits is a map of AI model names and the maximum number of tokens for them.
	TokenLimits = map[Model]uint{
		ModelGPT35Turbo:            4096,
		ModelGPT35TurboK16:         16385,
		ModelGPT35TurboInstruction: 4096,
		ModelGPT4:                  8192,
		ModelGPT4K32:               32768,
		ModelGPT4Preview:           4096, // but total input+output is 128000
		ModelGPT4VisionPreview:     4096, // but total input+output is 128000
		ModelGPT4o:                 4096, // but total input+output is 128000
	}
)

// Message is a struct of user message.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// Choice is a struct of response choice.
type Choice struct {
	Index        int          `json:"index"`
	Message      Message      `json:"message"`
	FinishReason FinishReason `json:"finish_reason"`
}

// Usage is additional information about the response limit usage.
type Usage struct {
	PromptTokens     uint `json:"prompt_tokens"`
	CompletionTokens uint `json:"completion_tokens"`
	TotalTokens      uint `json:"total_tokens"`
}

// CompletionRequest is a struct of request.
type CompletionRequest struct {
	Model    Model     `json:"model"`
	Messages []Message `json:"messages"`
	// optional
	MaxTokens        uint                `json:"max_tokens,omitempty"`
	User             string              `json:"user,omitempty"`
	Temperature      *float32            `json:"temperature,omitempty"`
	TopP             *float32            `json:"top_p,omitempty"`
	N                *uint               `json:"n,omitempty"`
	Stream           *bool               `json:"stream,omitempty"`
	Stop             *[]string           `json:"stop,omitempty"`
	PresencePenalty  *float32            `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float32            `json:"frequency_penalty,omitempty"`
	LogitBias        *map[string]float32 `json:"logit_bias,omitempty"`
}

func (c *CompletionRequest) marshal() (io.Reader, error) {
	if c.Model == "" {
		return nil, errors.Join(ErrRequiredParam, fmt.Errorf("model must not be empty"))
	}

	if len(c.Messages) == 0 {
		return nil, errors.Join(ErrRequiredParam, fmt.Errorf("messages must not be empty"))
	}

	if c.MaxTokens > 0 && c.MaxTokens > TokenLimits[c.Model] {
		return nil, errors.Join(
			ErrRequiredParam,
			fmt.Errorf("max tokens limit is %d, but gotten %d", TokenLimits[c.Model], c.MaxTokens),
		)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return bytes.NewReader(data), nil
}

func (c *CompletionRequest) build(ctx context.Context, auth *Params) (*http.Request, error) {
	body, err := c.marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, auth.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.Bearer))

	if auth.Organization != "" {
		req.Header.Set("OpenAI-Organization", auth.Organization)
	}

	return req, nil
}

// CompletionResponse is a struct of response.
type CompletionResponse struct {
	ID         string    `json:"id"`
	Object     string    `json:"object"`
	Created    int64     `json:"created"`
	Choices    []Choice  `json:"choices"`
	Usage      Usage     `json:"usage"`
	CreatedTs  time.Time `json:"-"`
	stopMarker string
}

func (r *CompletionResponse) build(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(r.Choices) == 0 {
		return errors.Join(ErrResponse, fmt.Errorf("empty response"))
	}

	r.CreatedTs = time.Unix(r.Created, 0)
	return nil
}

// String returns the first message of the response.
func (r *CompletionResponse) String() string {
	var (
		builder   strings.Builder
		hasMarker = r.stopMarker != ""
	)

	for i := range r.Choices {
		builder.WriteString(r.Choices[i].Message.Content)

		if hasMarker && (r.Choices[i].FinishReason == FinishReasonLength) {
			builder.WriteString(r.stopMarker)
		}
	}

	return builder.String()
}

// UsageInfo returns API tokens usage information.
func (r *CompletionResponse) UsageInfo() string {
	return fmt.Sprintf("prompt tokens: %d, completion tokens: %d, total tokens: %d",
		r.Usage.PromptTokens, r.Usage.CompletionTokens, r.Usage.TotalTokens,
	)
}

// Completion sends a request to the API and returns a response.
func Completion(ctx context.Context, client *http.Client, r *CompletionRequest, p Params) (*CompletionResponse, error) {
	body, err := commonRequest(ctx, client, r, p)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = body.Close()
	}()

	response := &CompletionResponse{stopMarker: p.StopMarker}
	if err = response.build(body); err != nil {
		return nil, err
	}

	return response, nil
}
