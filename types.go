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
	ModelDalle2           Model = "dall-e-2" // only for image requests
	ModelDalle3           Model = "dall-e-3" // only for image requests
	ModelGPT35Turbo       Model = "gpt-3.5-turbo"
	ModelGPT4             Model = "gpt-4"
	ModelGPT4Turbo        Model = "gpt-4-turbo"
	ModelGPT4o            Model = "gpt-4o"
	ModelGPT4oTurbo       Model = "gpt-4o-turbo"
	ModelGPT4oMini        Model = "gpt-4o-mini"
	ModelGPT41            Model = "gpt-4.1"      // only for image requests
	ModelGPT41Mini        Model = "gpt-4.1-mini" // only for image requests
	ModelGPT41Nano        Model = "gpt-4.1-nano" // only for image requests
	ModelGPT45Preview     Model = "gpt-4.5-preview"
	ModelGPTo1            Model = "o1"
	ModelGPTo1Mini        Model = "o1-mini"
	ModelGPTo1Preview     Model = "o1-preview"
	ModelGPTo1Pro         Model = "o1-pro"
	ModelGPTo3Mini        Model = "o3-mini"
	ModelDeepSeekChat     Model = "deepseek-chat"     // DeepSeek base model
	ModelDeepSeekReasoner Model = "deepseek-reasoner" // DeepSeek model with reasoning
)

// all models for image generation
var imageModels = map[Model]struct{}{ModelDalle2: {}, ModelDalle3: {}}

// MarshalJSON implements the json.Marshaler interface.
func (m *Model) MarshalJSON() ([]byte, error) {
	return marshalJSON(
		m,
		ModelDalle2, ModelDalle3,
		ModelGPT35Turbo, ModelGPT4, ModelGPT4Turbo, ModelGPT4o, ModelGPT4oTurbo, ModelGPT4oMini, ModelGPT45Preview,
		ModelGPTo1, ModelGPTo1Pro, ModelGPTo1Mini, ModelGPTo1Preview, ModelGPTo3Mini,
		ModelGPT41, ModelGPT41Mini, ModelGPT41Nano,
		ModelDeepSeekChat, ModelDeepSeekReasoner,
	)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Model) UnmarshalJSON(b []byte) error {
	return unMarshalJSON(
		m, b,
		ModelDalle2, ModelDalle3,
		ModelGPT35Turbo, ModelGPT4, ModelGPT4Turbo, ModelGPT4o, ModelGPT4oTurbo, ModelGPT4oMini, ModelGPT45Preview,
		ModelGPTo1, ModelGPTo1Pro, ModelGPTo1Mini, ModelGPTo1Preview, ModelGPTo3Mini,
		ModelGPT41, ModelGPT41Mini, ModelGPT41Nano,
		ModelDeepSeekChat, ModelDeepSeekReasoner,
	)
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
	var v = *t

	for _, value := range values {
		if v == value {
			return json.Marshal(string(v))
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
