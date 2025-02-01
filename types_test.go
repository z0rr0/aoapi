package aoapi

import (
	"errors"
	"testing"
)

func TestRole_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		role     Role
		expected string
		err      error
	}{
		{
			name:     "system",
			role:     RoleSystem,
			expected: `"system"`,
		},
		{
			name:     "user",
			role:     RoleUser,
			expected: `"user"`,
		},
		{
			name:     "assistant",
			role:     RoleAssistant,
			expected: `"assistant"`,
		},
		{
			name: "unknown",
			role: Role("unknown"),
			err:  ErrMarshalJSON,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.role.MarshalJSON()
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				return
			}

			if tc.err != nil {
				t.Fatalf("expected error, but got nil")
			}

			if s := string(data); s != tc.expected {
				t.Fatalf("expected: %q, got: %q", tc.expected, s)
			}
		})
	}
}

func TestRole_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		expected Role
		err      error
	}{
		{
			name:     "system",
			data:     `"system"`,
			expected: RoleSystem,
		},
		{
			name:     "user",
			data:     `"user"`,
			expected: RoleUser,
		},
		{
			name:     "assistant",
			data:     `"assistant"`,
			expected: RoleAssistant,
		},
		{
			name: "unknown",
			data: `"unknown"`,
			err:  ErrUnmarshalJSON,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			var role Role
			err := role.UnmarshalJSON([]byte(tc.data))
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				return
			}

			if tc.err != nil {
				t.Fatalf("expected error, but got nil")
			}

			if role != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, role)
			}
		})
	}
}

func TestModel_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		model    Model
		expected string
		err      error
	}{
		{
			name:     "gpt-3.5-turbo",
			model:    ModelGPT35Turbo,
			expected: `"gpt-3.5-turbo"`,
		},
		{
			name:     "gpt-4",
			model:    ModelGPT4,
			expected: `"gpt-4"`,
		},
		{
			name:     "gpt-4o",
			model:    ModelGPT4o,
			expected: `"gpt-4o"`,
		},
		{
			name:     "gpt-4o-mini",
			model:    ModelGPT4oMini,
			expected: `"gpt-4o-mini"`,
		},
		{
			name:     "gpt-4o-turbo",
			model:    ModelGPT4oTurbo,
			expected: `"gpt-4o-turbo"`,
		},
		{
			name:     "o1",
			model:    ModelGPTo1,
			expected: `"o1"`,
		},
		{
			name:     "o1-mini",
			model:    ModelGPTo1Mini,
			expected: `"o1-mini"`,
		},
		{
			name:     "o1-preview",
			model:    ModelGPTo1Preview,
			expected: `"o1-preview"`,
		},
		{
			name:     "o3-mini",
			model:    ModelGPTo3Mini,
			expected: `"o3-mini"`,
		},
		{
			name:     "deepseek-chat",
			model:    ModelDeepSeekChat,
			expected: `"deepseek-chat"`,
		},
		{
			name:     "deepseek-reasoner",
			model:    ModelDeepSeekReasoner,
			expected: `"deepseek-reasoner"`,
		},
		{
			name:     "dall-e-2",
			model:    ModelDalle2,
			expected: `"dall-e-2"`,
		},
		{
			name:     "dall-e-3",
			model:    ModelDalle3,
			expected: `"dall-e-3"`,
		},
		{
			name:  "unknown",
			model: Model("unknown"),
			err:   ErrMarshalJSON,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.model.MarshalJSON()
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				return
			}

			if tc.err != nil {
				t.Fatalf("expected error, but got nil")
			}

			if s := string(data); s != tc.expected {
				t.Fatalf("expected: %q, got: %q", tc.expected, s)
			}
		})
	}
}

func TestModel_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		expected Model
		err      error
	}{
		{
			name:     "gpt-3.5-turbo",
			data:     `"gpt-3.5-turbo"`,
			expected: ModelGPT35Turbo,
		},
		{
			name:     "gpt-4",
			data:     `"gpt-4"`,
			expected: ModelGPT4,
		},
		{
			name:     "gpt-4o",
			data:     `"gpt-4o"`,
			expected: ModelGPT4o,
		},
		{
			name:     "gpt-4o-mini",
			data:     `"gpt-4o-mini"`,
			expected: ModelGPT4oMini,
		},
		{
			name:     "gpt-4o-turbo",
			data:     `"gpt-4o-turbo"`,
			expected: ModelGPT4oTurbo,
		},
		{
			name:     "o1",
			data:     `"o1"`,
			expected: ModelGPTo1,
		},
		{
			name:     "o1-mini",
			data:     `"o1-mini"`,
			expected: ModelGPTo1Mini,
		},
		{
			name:     "o1-preview",
			data:     `"o1-preview"`,
			expected: ModelGPTo1Preview,
		},
		{
			name:     "o3-mini",
			data:     `"o3-mini"`,
			expected: ModelGPTo3Mini,
		},
		{
			name:     "deepseek-chan",
			data:     `"deepseek-chat"`,
			expected: ModelDeepSeekChat,
		},
		{
			name:     "deepseek-reasoner",
			data:     `"deepseek-reasoner"`,
			expected: ModelDeepSeekReasoner,
		},
		{
			name:     "dall-e-2",
			data:     `"dall-e-2"`,
			expected: ModelDalle2,
		},
		{
			name:     "dall-e-3",
			data:     `"dall-e-3"`,
			expected: ModelDalle3,
		},
		{
			name: "unknown",
			data: `"unknown"`,
			err:  ErrUnmarshalJSON,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			var model Model
			err := model.UnmarshalJSON([]byte(tc.data))
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				return
			}

			if tc.err != nil {
				t.Fatalf("expected error, but got nil")
			}

			if model != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, model)
			}

			_, isImage := imageModels[model]
			if _, ok := TokenLimits[model]; !(ok || isImage) {
				t.Errorf("model %v has no token limit", model)
			}
		})
	}
}

func TestFinishReason_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		reason   FinishReason
		expected string
		err      error
	}{
		{
			name:     "length",
			reason:   FinishReasonLength,
			expected: `"length"`,
		},
		{
			name:     "stop",
			reason:   FinishReasonStop,
			expected: `"stop"`,
		},
		{
			name:   "unknown",
			reason: FinishReason("unknown"),
			err:    ErrMarshalJSON,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.reason.MarshalJSON()
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				return
			}

			if tc.err != nil {
				t.Fatalf("expected error, but got nil")
			}

			if s := string(data); s != tc.expected {
				t.Fatalf("expected: %q, got: %q", tc.expected, s)
			}
		})
	}
}

func TestFinishReason_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		expected FinishReason
		err      error
	}{
		{
			name:     "length",
			data:     `"length"`,
			expected: FinishReasonLength,
		},
		{
			name:     "stop",
			data:     `"stop"`,
			expected: FinishReasonStop,
		},
		{
			name: "unknown",
			data: `"unknown"`,
			err:  ErrUnmarshalJSON,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			var reason FinishReason
			err := reason.UnmarshalJSON([]byte(tc.data))
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				return
			}

			if tc.err != nil {
				t.Fatalf("expected error, but got nil")
			}

			if reason != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, reason)
			}
		})
	}
}
