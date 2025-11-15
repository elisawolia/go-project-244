package formatters

import (
	"fmt"

	"code/internal/model"
)

func Format(tree []model.Node, format string) (string, error) {
	if format == "" {
		format = "stylish"
	}

	switch format {
	case "stylish":
		return formatStylish(tree), nil
	case "plain":
		return formatPlain(tree), nil
	case "json":
		return formatJSON(tree)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
