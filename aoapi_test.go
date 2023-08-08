package aoapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func compareResponses(a, b Response) bool {
	if a.ID != b.ID {
		return false
	}
	if a.Object != b.Object {
		return false
	}
	if a.Created != b.Created {
		return false
	}
	if a.Usage != b.Usage {
		return false
	}

	// compare Choices
	if len(a.Choices) != len(b.Choices) {
		return false
	}
	for i := range a.Choices {
		if a.Choices[i] != b.Choices[i] {
			return false
		}
	}

	return true
}

func TestRequestMarshal(t *testing.T) {
	var (
		temperature      float32 = 0.5
		topP             float32 = 0.6
		presencePenalty  float32 = 0.7
		frequencyPenalty float32 = 0.8
		n                uint    = 10
		stream                   = true
		logitBias                = map[string]float32{"a": 0.9}
	)

	testCases := []struct {
		name      string
		request   Request
		err       error
		errString string
		expected  []string
		length    int
	}{
		{
			name:      "empty request",
			request:   Request{},
			err:       RequiredParamError,
			errString: "required parameter is missing\nmodel must not be empty",
		},
		{
			name:      "empty messages",
			request:   Request{Model: ModelGPT4K32},
			err:       RequiredParamError,
			errString: "required parameter is missing\nmessages must not be empty",
		},
		{
			name:      "empty model",
			request:   Request{Messages: []Message{{Role: RoleSystem, Content: "Hello, world!"}}},
			err:       RequiredParamError,
			errString: "required parameter is missing\nmodel must not be empty",
		},
		{
			name: "valid request",
			request: Request{
				Model: ModelGPT35Turbo,
				Messages: []Message{
					{Role: RoleSystem, Content: "This is a system message"},
					{Role: RoleUser, Content: "This is a user message"},
				},
			},
			expected: []string{
				`"model":"gpt-3.5-turbo"`,
				`"messages":[`,
				`"role":"system"`,
				`"content":"This is a system message"`,
				`"role":"user"`,
				`"content":"This is a user message"`,
			},
			length: 144,
		},
		{
			name: "with name",
			request: Request{
				Model: ModelGPT35TurboK16,
				Messages: []Message{
					{Role: RoleAssistant, Content: "This is an assistant message", Name: "assistant"},
				},
			},
			expected: []string{
				`"model":"gpt-3.5-turbo-16k"`,
				`"messages":[`,
				`"role":"assistant"`,
				`"content":"This is an assistant message"`,
				`"name":"assistant"`,
			},
			length: 123,
		},
		{
			name: "with optional",
			request: Request{
				Model: ModelGPT4,
				Messages: []Message{
					{Role: RoleSystem, Content: "This is a system message"},
					{Role: RoleUser, Content: "This is a user message"},
				},
				// optional
				MaxTokens:        100,
				User:             "test",
				Temperature:      &temperature,
				TopP:             &topP,
				N:                &n,
				Stop:             &[]string{"\n\n"},
				Stream:           &stream,
				PresencePenalty:  &presencePenalty,
				FrequencyPenalty: &frequencyPenalty,
				LogitBias:        &logitBias,
			},
			expected: []string{
				`"model":"gpt-4"`,
				`"messages":[`,
				`"role":"system"`,
				`"content":"This is a system message"`,
				`"role":"user"`,
				`"content":"This is a user message"`,
				`"max_tokens":100`,
				`"user":"test"`,
				`"temperature":0.5`,
				`"top_p":0.6`,
				`"n":10`,
				`"stop":["\n\n"]`,
				`"stream":true`,
				`"presence_penalty":0.7`,
				`"frequency_penalty":0.8`,
				`"logit_bias":{"a":0.9}`,
			},
			length: 304,
		},
		{
			name: "partial optional",
			request: Request{
				Model: ModelGPT35Turbo,
				Messages: []Message{
					{Role: RoleSystem, Content: "This is a system message"},
					{Role: RoleUser, Content: "This is a user message"},
					{Role: RoleAssistant, Content: "This is an assistant message", Name: "assistant"},
				},
				MaxTokens:   250,
				Temperature: &temperature,
			},
			expected: []string{
				`"model":"gpt-3.5-turbo"`,
				`"messages":[`,
				`"role":"system"`,
				`"content":"This is a system message"`,
				`"role":"user"`,
				`"content":"This is a user message"`,
				`"role":"assistant"`,
				`"content":"This is an assistant message"`,
				`"name":"assistant"`,
				`"max_tokens":250`,
				`"temperature":0.5`,
			},
			length: 260,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			reader, err := tc.request.marshal()
			if err != nil {
				if errString := err.Error(); errString != tc.errString {
					t.Fatalf("unexpected error: %v", errString)
				}
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error type %v, but got %v", tc.err, err)
				}
				return
			}

			if tc.errString != "" {
				t.Fatalf("expected error %q", tc.errString)
			}

			buf := new(strings.Builder)
			_, err = io.Copy(buf, reader)
			if err != nil {
				t.Fatalf("failed to read request: %v", err)
			}

			s := buf.String()
			t.Logf("request: %s", s)

			if len(s) != tc.length {
				t.Fatalf("expected %d length, got %d", tc.length, len(s))
			}

			for _, expected := range tc.expected {
				if !strings.Contains(s, expected) {
					t.Fatalf("expected %q to contain %q", buf.String(), expected)
				}
			}
		})
	}
}

func TestCompletion(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{"id":"test","object":"chat.completion","created":1677652288,` +
			`"choices":[{"index":0,"message":{"content":"Message","role":"assistant"},"finish_reason":"stop"}],` +
			`"usage":{"prompt_tokens":4,"completion_tokens":6,"total_tokens":10}}`

		if _, err := fmt.Fprint(w, response); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &Request{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	response, err := Completion(context.Background(), client, request, s.URL, "test")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Response{
		ID:      "test",
		Object:  "chat.completion",
		Created: 1677652288,
		Choices: []Choice{{Index: 0, Message: Message{Content: "Message", Role: RoleAssistant}, FinishReason: "stop"}},
		Usage: Usage{
			PromptTokens:     4,
			CompletionTokens: 6,
			TotalTokens:      10,
		},
	}

	if r := *response; !compareResponses(expected, r) {
		t.Fatalf("expected %v, got %v", expected, r)
	}
}

func TestCompletionFailedStatus(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer s.Close()

	client := s.Client()
	request := &Request{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, s.URL, "test")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompletionFailedJSON(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err := fmt.Fprint(w, `{"id":"test","object":"chat.completion,`); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &Request{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, s.URL, "test")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompletionFailedURL(t *testing.T) {
	client := http.DefaultClient
	request := &Request{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, "http://127.0.0.1:99999", "test")

	if err == nil {
		t.Fatal("expected error")
	}
}
