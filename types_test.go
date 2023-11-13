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
			name:     "gpt-3.5-turbo-16k",
			model:    ModelGPT35TurboK16,
			expected: `"gpt-3.5-turbo-16k"`,
		},
		{
			name:     "gpt-3.5-turbo-instruct",
			model:    ModelGPT35TurboInstruction,
			expected: `"gpt-3.5-turbo-instruct"`,
		},
		{
			name:     "gpt-4",
			model:    ModelGPT4,
			expected: `"gpt-4"`,
		},
		{
			name:     "gpt-4-32k",
			model:    ModelGPT4K32,
			expected: `"gpt-4-32k"`,
		},
		{
			name:     "gpt-4-1106-preview",
			model:    ModelGPT4Preview,
			expected: `"gpt-4-1106-preview"`,
		},
		{
			name:     "gpt-4-1106-vision-preview",
			model:    ModelGPT4VisionPreview,
			expected: `"gpt-4-1106-vision-preview"`,
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
			name:     "gpt-3.5-turbo-16k",
			data:     `"gpt-3.5-turbo-16k"`,
			expected: ModelGPT35TurboK16,
		},
		{
			name:     "gpt-3.5-turbo-instruct",
			data:     `"gpt-3.5-turbo-instruct"`,
			expected: ModelGPT35TurboInstruction,
		},
		{
			name:     "gpt-4",
			data:     `"gpt-4"`,
			expected: ModelGPT4,
		},
		{
			name:     "gpt-4-32k",
			data:     `"gpt-4-32k"`,
			expected: ModelGPT4K32,
		},
		{
			name:     "gpt-4-1106-preview",
			data:     `"gpt-4-1106-preview"`,
			expected: ModelGPT4Preview,
		},
		{
			name:     "gpt-4-1106-vision-preview",
			data:     `"gpt-4-1106-vision-preview"`,
			expected: ModelGPT4VisionPreview,
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
