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
	return marshalJSON(r, RoleSystem, RoleUser, RoleAssistant)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Role) UnmarshalJSON(b []byte) error {
	return unMarshalJSON(r, b, RoleSystem, RoleUser, RoleAssistant)
}

// Model is a type of AI model name.
type Model string

// AI model names.
const (
	ModelGPT35Turbo            Model = "gpt-3.5-turbo"
	ModelGPT35TurboK16         Model = "gpt-3.5-turbo-16k"
	ModelGPT35TurboInstruction Model = "gpt-3.5-turbo-instruct"
	ModelGPT4                  Model = "gpt-4"
	ModelGPT4K32               Model = "gpt-4-32k"
)

// MarshalJSON implements the json.Marshaler interface.
func (m *Model) MarshalJSON() ([]byte, error) {
	return marshalJSON(m, ModelGPT35Turbo, ModelGPT35TurboK16, ModelGPT35TurboInstruction, ModelGPT4, ModelGPT4K32)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Model) UnmarshalJSON(b []byte) error {
	return unMarshalJSON(m, b, ModelGPT35Turbo, ModelGPT35TurboK16, ModelGPT35TurboInstruction, ModelGPT4, ModelGPT4K32)
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
	return marshalJSON(f, FinishReasonLength, FinishReasonStop)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *FinishReason) UnmarshalJSON(b []byte) error {
	return unMarshalJSON(f, b, FinishReasonLength, FinishReasonStop)
}

// StringCommonType is a generic interface for custom string based types.
type StringCommonType interface {
	Role | Model | FinishReason
}

// marshalJSON is a generic function for custom types JSON marshal.
func marshalJSON[T StringCommonType](t *T, values ...T) ([]byte, error) {
	var (
		v = *t
		s = string(v)
	)

	for _, value := range values {
		if v == value {
			return json.Marshal(s)
		}
	}

	return nil, errors.Join(ErrMarshalJSON, fmt.Errorf("invalid value: %v", v))
}

// unMarshalJSON is a generic function for custom types JSON unmarshal.
func unMarshalJSON[T StringCommonType](t *T, b []byte, values ...T) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return errors.Join(ErrUnmarshalJSON, fmt.Errorf("invalid string value: %v", string(b)))
	}

	v := T(s)

	for _, value := range values {
		if v == value {
			*t = v
			return nil
		}
	}

	return errors.Join(ErrUnmarshalJSON, fmt.Errorf("invalid value: %v", v))
}
