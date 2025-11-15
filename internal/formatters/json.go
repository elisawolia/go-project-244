package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/model"
)

func formatJSON(tree []model.Node) (string, error) {
	data, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return "", fmt.Errorf("cannot marshal diff to json: %w", err)
	}
	return string(data), nil
}
