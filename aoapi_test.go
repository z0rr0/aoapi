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

func compareCompletionResponses(a, b CompletionResponse) bool {
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

func TestCompletionRequestMarshal(t *testing.T) {
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
		request   CompletionRequest
		err       error
		errString string
		expected  []string
		length    int
	}{
		{
			name:      "empty request",
			request:   CompletionRequest{},
			err:       ErrRequiredParam,
			errString: "required parameter is missing\nmodel must not be empty",
		},
		{
			name:      "empty messages",
			request:   CompletionRequest{Model: ModelGPT4K32},
			err:       ErrRequiredParam,
			errString: "required parameter is missing\nmessages must not be empty",
		},
		{
			name:      "empty model",
			request:   CompletionRequest{Messages: []Message{{Role: RoleSystem, Content: "Hello, world!"}}},
			err:       ErrRequiredParam,
			errString: "required parameter is missing\nmodel must not be empty",
		},
		{
			name: "valid request",
			request: CompletionRequest{
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
			request: CompletionRequest{
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
			request: CompletionRequest{
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
			request: CompletionRequest{
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
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("failed content type header: %q", ct)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test" {
			t.Errorf("failed authorization header: %q", auth)
		}
		if org := r.Header.Get("OpenAI-Organization"); org != "test-org" {
			t.Errorf("failed organization header: %q", org)
		}

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
	request := &CompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}

	params := Params{Bearer: "test", URL: s.URL, Organization: "test-org"}
	response, err := Completion(context.Background(), client, request, params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := CompletionResponse{
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

	if r := *response; !compareCompletionResponses(expected, r) {
		t.Fatalf("expected %v, got %v", expected, r)
	}
}

func TestCompletionFailedStatus(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test", http.StatusBadGateway)
	}))
	defer s.Close()

	client := s.Client()
	request := &CompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompletionFailedWithError(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")

		response := `{"error":{"message":"This model's maximum context length is 4097 tokens.",` +
			`"type": "invalid_request_error","param": "messages","code": "context_length_exceeded"}}`

		if _, err := fmt.Fprint(w, response); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &CompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrResponse) {
		t.Fatalf("expected %v, got %v", ErrResponse, err)
	}

	expected := "failed response\nstatus code 400\n"
	expected += `type="invalid_request_error", param="messages", code="context_length_exceeded": `
	expected += `This model's maximum context length is 4097 tokens.`

	if e := err.Error(); e != expected {
		t.Fatalf("expected %q, got %q", expected, e)
	}
}

func TestCompletionFailedRequest(t *testing.T) {
	client := http.DefaultClient
	request := &CompletionRequest{Model: ModelGPT35Turbo} // no messages
	_, err := Completion(context.Background(), client, request, Params{Bearer: "test", URL: ":"})

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrRequiredParam) {
		t.Fatalf("expected %v, got %v", ErrResponse, err)
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
	request := &CompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}

	if e := err.Error(); !strings.HasPrefix(e, "failed to unmarshal response") {
		t.Fatalf("expected %q, got %q", "failed to unmarshal response", e)
	}
}

func TestCompletionFailedURL(t *testing.T) {
	client := http.DefaultClient
	request := &CompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, Params{Bearer: "test", URL: "http://127.0.0.1:99999"})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompletionFailedContent(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"id":"test","object":"chat.completion","created":1677652288,` +
			`"usage":{"prompt_tokens":4,"completion_tokens":6,"total_tokens":10}}`

		if _, err := fmt.Fprint(w, response); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &CompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []Message{
			{Role: RoleSystem, Content: "This is a system message"},
			{Role: RoleUser, Content: "This is a user message"},
		},
		MaxTokens: 100,
	}
	_, err := Completion(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrResponse) {
		t.Fatalf("expected %v, got %v", ErrResponse, err)
	}
}

func TestCompletionResponse_String(t *testing.T) {
	testCases := []struct {
		name     string
		response CompletionResponse
		expected string
	}{
		{
			name: "empty",
			response: CompletionResponse{
				ID:      "test",
				Object:  "chat.completion",
				Created: 1677652288,
			},
		},
		{
			name: "message",
			response: CompletionResponse{
				ID:      "test",
				Object:  "chat.completion",
				Created: 1677652288,
				Choices: []Choice{
					{
						Index:        0,
						Message:      Message{Content: "This is ", Role: RoleAssistant},
						FinishReason: FinishReasonStop,
					},
					{
						Index:        1,
						Message:      Message{Content: "a message.", Role: RoleAssistant},
						FinishReason: FinishReasonStop,
					},
				},
			},
			expected: "This is a message.",
		},
		{
			name: "length",
			response: CompletionResponse{
				ID:      "test",
				Object:  "chat.completion",
				Created: 1677652288,
				Choices: []Choice{
					{
						Index:        0,
						Message:      Message{Content: "This is a message.", Role: RoleAssistant},
						FinishReason: FinishReasonLength,
					},
				},
				stopMarker: " [reason=length]",
			},
			expected: "This is a message. [reason=length]",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			if s := tc.response.String(); s != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, s)
			}
		})
	}
}

func TestCompletionResponse_UsageInfo(t *testing.T) {
	response := CompletionResponse{
		ID:      "test",
		Object:  "chat.completion",
		Created: 1677652288,
		Choices: []Choice{
			{Message: Message{Content: "This is a message", Role: RoleAssistant}, FinishReason: FinishReasonStop},
		},
		Usage: Usage{
			PromptTokens:     4,
			CompletionTokens: 6,
			TotalTokens:      10,
		},
	}

	expected := "prompt tokens: 4, completion tokens: 6, total tokens: 10"

	if s := response.UsageInfo(); s != expected {
		t.Errorf("expected %v, got %v", expected, s)
	}
}

func TestCompletionRequestMaxTokens(t *testing.T) {
	testCases := []struct {
		name      string
		model     Model
		maxTokens uint
		withErr   bool
	}{
		{
			name:  "valid_no_limit",
			model: ModelGPT35Turbo,
		},
		{
			name:      "valid_with_limit",
			model:     ModelGPT35TurboK16,
			maxTokens: TokenLimits[ModelGPT35TurboK16] - 1,
		},
		{
			name:      "failed",
			model:     ModelGPT4,
			maxTokens: TokenLimits[ModelGPT4] + 1,
			withErr:   true,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			request := CompletionRequest{
				Model:     tc.model,
				Messages:  []Message{{Role: RoleSystem, Content: "Hello, world!"}},
				MaxTokens: tc.maxTokens,
			}
			_, err := request.marshal()

			switch {
			case err != nil && !tc.withErr:
				t.Errorf("unexpected error: %v", err)
			case err == nil && tc.withErr:
				t.Errorf("expected error")
			case err != nil:
				if !errors.Is(err, ErrRequiredParam) {
					t.Errorf("expected %v, got %v", ErrRequiredParam, err)
				}

				if s := err.Error(); !strings.Contains(s, "max tokens limit") {
					t.Errorf("not found required error string, got %q", s)
				}
			}
		})
	}
}
