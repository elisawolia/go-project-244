package formatters

import (
	"fmt"
	"strings"

	"code/internal/model"
	"code/internal/utils"
)

func formatPlain(tree []model.Node) string {
	lines := buildPlainLines(tree, "")
	return strings.Join(lines, "\n")
}

func buildPlainLines(nodes []model.Node, parentPath string) []string {
	var result []string

	for _, n := range nodes {
		fullPath := n.Key
		if parentPath != "" {
			fullPath = parentPath + "." + n.Key
		}

		switch n.Type {
		case model.NodeNested:
			result = append(result, buildPlainLines(n.Children, fullPath)...)

		case model.NodeAdded:
			value := plainFormatValue(n.Value)
			line := fmt.Sprintf("Property '%s' was added with value: %s", fullPath, value)
			result = append(result, line)

		case model.NodeRemoved:
			line := fmt.Sprintf("Property '%s' was removed", fullPath)
			result = append(result, line)

		case model.NodeUpdated:
			oldVal := plainFormatValue(n.OldValue)
			newVal := plainFormatValue(n.NewValue)
			line := fmt.Sprintf("Property '%s' was updated. From %s to %s", fullPath, oldVal, newVal)
			result = append(result, line)

		case model.NodeUnchanged:
		}
	}

	return result
}

func plainFormatValue(v interface{}) string {
	if _, ok := utils.ToComplexStruct(v); ok {
		return "[complex value]"
	}

	switch val := v.(type) {
	case []interface{}:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
