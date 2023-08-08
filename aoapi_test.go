package aoapi

import (
	"io"
	"strings"
	"testing"
)

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
		name     string
		request  Request
		err      string
		expected []string
		length   int
	}{
		{
			name:    "empty request",
			request: Request{},
			err:     "messages must not be empty",
		},
		{
			name:    "empty messages",
			request: Request{Model: ModelGPT4K32},
			err:     "messages must not be empty",
		},
		{
			name:    "empty model",
			request: Request{Messages: []Message{{Role: RoleSystem, Content: "Hello, world!"}}},
			err:     "model must not be empty",
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
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			reader, err := tc.request.marshal()
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if tc.err != "" {
				t.Fatalf("expected error %q", tc.err)
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
