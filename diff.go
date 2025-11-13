package code

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"code/internal/parser"
)

type nodeType string

const (
	nodeAdded     nodeType = "added"
	nodeRemoved   nodeType = "removed"
	nodeUnchanged nodeType = "unchanged"
	nodeUpdated   nodeType = "updated"
	nodeNested    nodeType = "nested"
)

type node struct {
	Key      string
	Type     nodeType
	Value    interface{} // для added/removed/unchanged
	OldValue interface{} // для updated
	NewValue interface{} // для updated
	Children []node      // для nested
}

func GenDiff(path1, path2 string) (string, error) {
	data1, err := parser.Parse(path1)
	if err != nil {
		return "", err
	}

	data2, err := parser.Parse(path2)
	if err != nil {
		return "", err
	}

	tree := buildAST(data1, data2)

	return formatStylish(tree), nil
}

func buildAST(m1, m2 map[string]interface{}) []node {
	keysSet := make(map[string]struct{})
	for k := range m1 {
		keysSet[k] = struct{}{}
	}
	for k := range m2 {
		keysSet[k] = struct{}{}
	}

	keys := make([]string, 0, len(keysSet))
	for k := range keysSet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]node, 0, len(keys))

	for _, key := range keys {
		v1, ok1 := m1[key]
		v2, ok2 := m2[key]

		switch {
		case ok1 && ok2:
			m1Child, isMap1 := toMap(v1)
			m2Child, isMap2 := toMap(v2)

			if isMap1 && isMap2 {
				child := node{
					Key:      key,
					Type:     nodeNested,
					Children: buildAST(m1Child, m2Child),
				}
				result = append(result, child)
				continue
			}

			if reflect.DeepEqual(v1, v2) {
				result = append(result, node{
					Key:   key,
					Type:  nodeUnchanged,
					Value: v1,
				})
			} else {
				result = append(result, node{
					Key:      key,
					Type:     nodeUpdated,
					OldValue: v1,
					NewValue: v2,
				})
			}

		case ok1 && !ok2:
			result = append(result, node{
				Key:   key,
				Type:  nodeRemoved,
				Value: v1,
			})

		case !ok1 && ok2:
			result = append(result, node{
				Key:   key,
				Type:  nodeAdded,
				Value: v2,
			})
		}
	}

	return result
}

func toMap(v interface{}) (map[string]interface{}, bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{}, len(m))
		for k, val := range m {
			if nested, ok := toMap(val); ok {
				res[k] = nested
			} else {
				res[k] = val
			}
		}
		return res, true
	case map[interface{}]interface{}:
		res := make(map[string]interface{}, len(m))
		for k, val := range m {
			keyStr, ok := k.(string)
			if !ok {
				keyStr = fmt.Sprint(k)
			}
			if nested, ok := toMap(val); ok {
				res[keyStr] = nested
			} else {
				res[keyStr] = val
			}
		}
		return res, true
	default:
		return nil, false
	}
}

func formatStylish(tree []node) string {
	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString(formatNodes(tree, 1))
	b.WriteString("}")
	return b.String()
}

func formatNodes(nodes []node, depth int) string {
	var b strings.Builder

	for _, n := range nodes {
		switch n.Type {
		case nodeNested:
			indent := makeIndent(depth, ' ')
			fmt.Fprintf(&b, "%s%s: {\n", indent, n.Key)
			b.WriteString(formatNodes(n.Children, depth+1))
			fmt.Fprintf(&b, "%s}\n", strings.Repeat(" ", depth*4))

		case nodeUnchanged:
			indent := makeIndent(depth, ' ')
			fmt.Fprintf(&b, "%s%s: %s\n", indent, n.Key, formatValue(n.Value, depth))

		case nodeAdded:
			indent := makeIndent(depth, '+')
			fmt.Fprintf(&b, "%s%s: %s\n", indent, n.Key, formatValue(n.Value, depth))

		case nodeRemoved:
			indent := makeIndent(depth, '-')
			fmt.Fprintf(&b, "%s%s: %s\n", indent, n.Key, formatValue(n.Value, depth))

		case nodeUpdated:
			indentOld := makeIndent(depth, '-')
			indentNew := makeIndent(depth, '+')
			fmt.Fprintf(&b, "%s%s: %s\n", indentOld, n.Key, formatValue(n.OldValue, depth))
			fmt.Fprintf(&b, "%s%s: %s\n", indentNew, n.Key, formatValue(n.NewValue, depth))
		}
	}

	return b.String()
}

func makeIndent(depth int, sign rune) string {
	baseIndent := depth*4 - 2
	if baseIndent < 0 {
		baseIndent = 0
	}
	return strings.Repeat(" ", baseIndent) + string(sign) + " "
}

func formatValue(v interface{}, depth int) string {
	if m, ok := toMap(v); ok {
		return formatMap(m, depth+1)
	}

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
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
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
		indent := strings.Repeat(" ", depth*4)
		fmt.Fprintf(&b, "%s%s: %s\n", indent, key, formatValue(val, depth))
	}
	fmt.Fprintf(&b, "%s}", strings.Repeat(" ", (depth-1)*4))

	return b.String()
}
