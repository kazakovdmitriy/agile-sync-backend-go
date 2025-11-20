package utils

import (
	"encoding/json"
	"fmt"
)

// MapToStruct безопасно конвертирует map[string]interface{} в struct
func MapToStruct(input map[string]interface{}, target interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input map: %w", err)
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal into target struct: %w", err)
	}
	return nil
}
