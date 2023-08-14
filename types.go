package aoapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrMarshalJSON is an error that occurs when a custom type is marshaled to JSON.
	ErrMarshalJSON = errors.New("failed JSON marshal")

	// ErrUnmarshalJSON is an error that occurs when a custom type is unmarshaled from JSON.
	ErrUnmarshalJSON = errors.New("failed JSON unmarshal")
)

// Role is a type of user message role.
type Role string

// User message roles.
const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// MarshalJSON implements the json.Marshaler interface.
func (r *Role) MarshalJSON() ([]byte, error) {
	return marshalJSON[Role](r, RoleSystem, RoleUser, RoleAssistant)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Role) UnmarshalJSON(b []byte) error {
	return unMarshalJSON[Role](r, b, RoleSystem, RoleUser, RoleAssistant)
}

// Model is a type of AI model name.
type Model string

// AI model names.
const (
	ModelGPT35Turbo    Model = "gpt-3.5-turbo"
	ModelGPT35TurboK16 Model = "gpt-3.5-turbo-16k"
	ModelGPT4          Model = "gpt-4"
	ModelGPT4K32       Model = "gpt-4-32k"
)

// MarshalJSON implements the json.Marshaler interface.
func (m *Model) MarshalJSON() ([]byte, error) {
	return marshalJSON[Model](m, ModelGPT35Turbo, ModelGPT35TurboK16, ModelGPT4, ModelGPT4K32)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Model) UnmarshalJSON(b []byte) error {
	return unMarshalJSON[Model](m, b, ModelGPT35Turbo, ModelGPT35TurboK16, ModelGPT4, ModelGPT4K32)
}

// FinishReason is a type of response finish reason.
type FinishReason string

// Finish reasons variants.
const (
	FinishReasonLength FinishReason = "length"
	FinishReasonStop   FinishReason = "stop"
)

// MarshalJSON implements the json.Marshaler interface.
func (f *FinishReason) MarshalJSON() ([]byte, error) {
	return marshalJSON[FinishReason](f, FinishReasonLength, FinishReasonStop)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *FinishReason) UnmarshalJSON(b []byte) error {
	return unMarshalJSON[FinishReason](f, b, FinishReasonLength, FinishReasonStop)
}

// marshalJSON is a generic function for custom types JSON marshal.
func marshalJSON[T Role | Model | FinishReason](t *T, values ...T) ([]byte, error) {
	v := *t

	for _, value := range values {
		if v == value {
			return json.Marshal(string(v))
		}
	}

	return nil, errors.Join(ErrMarshalJSON, fmt.Errorf("invalid value: %v", v))
}

// unMarshalJSON is a generic function for custom types JSON unmarshal.
func unMarshalJSON[T Role | Model | FinishReason](t *T, b []byte, values ...T) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return errors.Join(ErrUnmarshalJSON, fmt.Errorf("invalid string value: %v", str))
	}

	v := T(str)

	for _, value := range values {
		if v == value {
			*t = v
			return nil
		}
	}

	return errors.Join(ErrUnmarshalJSON, fmt.Errorf("invalid value: %v", v))
}
