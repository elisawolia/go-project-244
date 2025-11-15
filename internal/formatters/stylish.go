package formatters

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"code/internal/model"
	"code/internal/utils"
)

const indentSize = 4
const signOffset = 2
const printPattern = "%s%s: %s\n"

func formatStylish(tree []model.Node) string {
	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString(formatNodes(tree, 1))
	b.WriteString("}")
	return b.String()
}

func formatNodes(nodes []model.Node, depth int) string {
	var b strings.Builder

	for _, n := range nodes {
		switch n.Type {
		case model.NodeNested:
			indent := makeIndent(depth, ' ')
			fmt.Fprintf(&b, "%s%s: {\n", indent, n.Key)
			b.WriteString(formatNodes(n.Children, depth+1))
			fmt.Fprintf(&b, "%s}\n", strings.Repeat(" ", depth*indentSize))

		case model.NodeUnchanged:
			indent := makeIndent(depth, ' ')
			fmt.Fprintf(&b, printPattern, indent, n.Key, formatValue(n.Value, depth))

		case model.NodeAdded:
			indent := makeIndent(depth, '+')
			fmt.Fprintf(&b, printPattern, indent, n.Key, formatValue(n.Value, depth))

		case model.NodeRemoved:
			indent := makeIndent(depth, '-')
			fmt.Fprintf(&b, printPattern, indent, n.Key, formatValue(n.Value, depth))

		case model.NodeUpdated:
			indentOld := makeIndent(depth, '-')
			indentNew := makeIndent(depth, '+')
			fmt.Fprintf(&b, printPattern, indentOld, n.Key, formatValue(n.OldValue, depth))
			fmt.Fprintf(&b, printPattern, indentNew, n.Key, formatValue(n.NewValue, depth))
		}
	}

	return b.String()
}

func makeIndent(depth int, sign rune) string {
	baseIndent := depth*indentSize - signOffset
	if baseIndent < 0 {
		baseIndent = 0
	}
	return strings.Repeat(" ", baseIndent) + string(sign) + " "
}

func formatValue(v interface{}, depth int) string {
	if complexStr, ok := tryFormatComplex(v, depth); ok {
		return complexStr
	}

	return formatPrimitive(v)
}

func tryFormatComplex(v interface{}, depth int) (string, bool) {
	complexVal, ok := utils.ToComplexStruct(v)
	if !ok {
		return "", false
	}

	switch val := complexVal.(type) {
	case map[string]interface{}:
		return formatMap(val, depth+1), true
	case []interface{}:
		return formatSlice(val, depth+1), true
	default:
		return "", false
	}
}

func formatPrimitive(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		if val == math.Trunc(val) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func formatMap(m map[string]interface{}, depth int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("{\n")
	for _, key := range keys {
		val := m[key]
		indent := strings.Repeat(" ", depth*indentSize)
		fmt.Fprintf(&b, printPattern, indent, key, formatValue(val, depth))
	}
	fmt.Fprintf(&b, "%s}", strings.Repeat(" ", (depth-1)*indentSize))

	return b.String()
}

func formatSlice(s []interface{}, depth int) string {
	var b strings.Builder

	indent := strings.Repeat(" ", depth*indentSize)
	closingIndent := strings.Repeat(" ", (depth-1)*indentSize)

	b.WriteString("[\n")
	for i, item := range s {
		fmt.Fprintf(&b, "%s%s", indent, formatValue(item, depth))
		if i < len(s)-1 {
			b.WriteString("\n")
		}
	}
	if len(s) > 0 {
		b.WriteString("\n")
	}
	b.WriteString(closingIndent)
	b.WriteString("]")

	return b.String()
}
