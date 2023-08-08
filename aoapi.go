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
	ErrRequiredParam = fmt.Errorf("required parameter is missing")

	// ErrResponse is an error that occurs when the response is empty.
	ErrResponse = fmt.Errorf("failed response")
)

// Role is a type of user message role.
type Role string

// User message roles.
const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Model is a type of AI model name.
type Model string

// AI model names.
const (
	ModelGPT35Turbo    Model = "gpt-3.5-turbo"
	ModelGPT35TurboK16 Model = "gpt-3.5-turbo-16k"
	ModelGPT4          Model = "gpt-4"
	ModelGPT4K32       Model = "gpt-4-32k"
)

// FinishReason is a type of response finish reason.
type FinishReason string

// Finish reasons variants.
const (
	FinishReasonLength FinishReason = "length"
	FinishReasonStop   FinishReason = "stop"
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
	if r.Model == "" {
		return nil, errors.Join(ErrRequiredParam, fmt.Errorf("model must not be empty"))
	}

	if len(r.Messages) == 0 {
		return nil, errors.Join(ErrRequiredParam, fmt.Errorf("messages must not be empty"))
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return bytes.NewReader(data), nil
}

func (r *Request) build(ctx context.Context, auth Auth) (*http.Request, error) {
	body, err := r.marshal()
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

// Response is a struct of response.
type Response struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Created   int64     `json:"created"`
	Choices   []Choice  `json:"choices"`
	Usage     Usage     `json:"usage"`
	CreatedTs time.Time `json:"-"`
}

func (response *Response) build(resp *http.Response) error {
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return errors.Join(ErrResponse, fmt.Errorf("empty response"))
	}

	response.CreatedTs = time.Unix(response.Created, 0)
	return nil
}

// String returns the first message of the response.
func (response *Response) String() string {
	const reasonMarker = " [reason=length]"
	var b strings.Builder

	for i := range response.Choices {
		b.WriteString(response.Choices[i].Message.Content)

		if response.Choices[i].FinishReason == FinishReasonLength {
			b.WriteString(reasonMarker)
		}
	}
	return b.String()
}

// UsageInfo returns API tokens usage information.
func (response *Response) UsageInfo() string {
	return fmt.Sprintf("prompt tokens: %d, completion tokens: %d, total tokens: %d",
		response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens,
	)
}

// Auth is a struct of API authentication information.
type Auth struct {
	Bearer       string
	Organization string
	URL          string
}

// Completion sends a request to the API and returns a response.
func Completion(ctx context.Context, client *http.Client, r *Request, auth Auth) (*Response, error) {
	request, err := r.build(ctx, auth)
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Join(ErrResponse, fmt.Errorf("status code %d", resp.StatusCode))
	}

	response := &Response{}
	if err = response.build(resp); err != nil {
		return nil, err
	}

	return response, nil
}
