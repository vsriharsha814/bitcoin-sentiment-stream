package config

import (
	"encoding/json"
	"fmt"
)

// QuestionMapping maps numeric keys → full question text.
var QuestionMapping map[string]string

// LoadQuestionMapping takes the raw JSON bytes of mapping.json.
func LoadQuestionMapping(raw []byte) error {
	if err := json.Unmarshal(raw, &QuestionMapping); err != nil {
		return fmt.Errorf("parse mapping: %w", err)
	}
	return nil
}
