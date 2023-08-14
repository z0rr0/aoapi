package aoapi

import (
	"encoding/json"
	"errors"
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
	switch role := *r; role {
	case RoleSystem, RoleUser, RoleAssistant:
		return json.Marshal(string(role))
	default:
		return nil, errors.New("invalid role")
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Role) UnmarshalJSON(b []byte) error {
	var roleStr string
	if err := json.Unmarshal(b, &roleStr); err != nil {
		return err
	}

	switch role := Role(roleStr); role {
	case RoleSystem, RoleUser, RoleAssistant:
		*r = role
	default:
		return errors.New("invalid role")
	}

	return nil
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
	switch model := *m; model {
	case ModelGPT35Turbo, ModelGPT35TurboK16, ModelGPT4, ModelGPT4K32:
		return json.Marshal(string(model))
	default:
		return nil, errors.New("invalid model")
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Model) UnmarshalJSON(b []byte) error {
	var modelStr string
	if err := json.Unmarshal(b, &modelStr); err != nil {
		return err
	}

	switch model := Model(modelStr); model {
	case ModelGPT35Turbo, ModelGPT35TurboK16, ModelGPT4, ModelGPT4K32:
		*m = model
	default:
		return errors.New("invalid model")
	}

	return nil
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
	switch reason := *f; reason {
	case FinishReasonLength, FinishReasonStop:
		return json.Marshal(string(reason))
	default:
		return nil, errors.New("invalid finish reason")
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *FinishReason) UnmarshalJSON(b []byte) error {
	var reasonStr string
	if err := json.Unmarshal(b, &reasonStr); err != nil {
		return err
	}

	switch reason := FinishReason(reasonStr); reason {
	case FinishReasonLength, FinishReasonStop:
		*f = reason
	default:
		return errors.New("invalid finish reason")
	}

	return nil
}
