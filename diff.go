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

func GenDiff(path1, path2 string) (string, error) {
	data1, err := parser.Parse(path1)
	if err != nil {
		return "", err
	}

	data2, err := parser.Parse(path2)
	if err != nil {
		return "", err
	}

	diff := buildFlatDiff(data1, data2)
	return diff, nil
}

func buildFlatDiff(m1, m2 map[string]interface{}) string {
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

	var b strings.Builder
	b.WriteString("{\n")

	for _, key := range keys {
		v1, ok1 := m1[key]
		v2, ok2 := m2[key]

		switch {
		case ok1 && ok2:
			if isEqual(v1, v2) {
				fmt.Fprintf(&b, "    %s: %s\n", key, formatValue(v1))
			} else {
				fmt.Fprintf(&b, "  - %s: %s\n", key, formatValue(v1))
				fmt.Fprintf(&b, "  + %s: %s\n", key, formatValue(v2))
			}
		case ok1 && !ok2:
			fmt.Fprintf(&b, "  - %s: %s\n", key, formatValue(v1))
		case !ok1 && ok2:
			fmt.Fprintf(&b, "  + %s: %s\n", key, formatValue(v2))
		}
	}

	b.WriteString("}")
	return b.String()
}

func isEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func formatValue(v interface{}) string {
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
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
