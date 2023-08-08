package aoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Role is a type of user message role.
type Role string

// Model is a type of AI model name.
type Model string

// User message roles.
const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// AI model names.
const (
	ModelGPT35Turbo    Model = "gpt-3.5-turbo"
	ModelGPT35TurboK16 Model = "gpt-3.5-turbo-16k"
	ModelGPT4          Model = "gpt-4"
	ModelGPT4K32       Model = "gpt-4-32k"
)

// Message is a struct of user message.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// Choice is a struct of response choice.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage is additional information about the response limit usage.
type Usage struct {
	PromptTokens     uint `json:"prompt_tokens"`
	CompletionTokens uint `json:"completion_tokens"`
	TotalTokens      uint `json:"total_tokens"`
}

// Response is a struct of response.
type Response struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Created   int64     `json:"created"`
	Choices   []Choice  `json:"choices"`
	Usage     Usage     `json:"usage"`
	CreatedTs time.Time `json:"-"`
}

// Request is a struct of request.
type Request struct {
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

func (r *Request) marshal() (io.Reader, error) {
	if len(r.Messages) == 0 {
		return nil, fmt.Errorf("messages must not be empty")
	}

	if r.Model == "" {
		return nil, fmt.Errorf("model must not be empty")
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return bytes.NewReader(data), nil
}

func (r *Request) build(ctx context.Context, uri, bearer string) (*http.Request, error) {
	body, err := r.marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))

	return req, nil
}

func (response *Response) build(resp *http.Response) error {
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	response.CreatedTs = time.Unix(response.Created, 0)
	return nil
}

// Completion sends a request to the API and returns a response.
func Completion(ctx context.Context, client *http.Client, r *Request, uri, bearer string) (*Response, error) {
	request, err := r.build(ctx, uri, bearer)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	response := &Response{}
	if err = response.build(resp); err != nil {
		return nil, err
	}

	return response, nil
}