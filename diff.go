package code

import (
	"reflect"
	"sort"

	"code/internal/formatters"
	"code/internal/model"
	"code/internal/parser"
	"code/internal/utils"
)

func GenDiff(path1, path2, format string) (string, error) {
	data1, err := parser.Parse(path1)
	if err != nil {
		return "", err
	}

	data2, err := parser.Parse(path2)
	if err != nil {
		return "", err
	}

	tree := buildAST(data1, data2)

	return formatters.Format(tree, format)
}

func buildAST(m1, m2 map[string]interface{}) []model.Node {
	keys := collectKeys(m1, m2)

	result := make([]model.Node, 0, len(keys))
	for _, key := range keys {
		result = append(result, buildNode(key, m1, m2))
	}

	return result
}

func collectKeys(m1, m2 map[string]interface{}) []string {
	keysSet := make(map[string]struct{}, len(m1)+len(m2))

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

	return keys
}

func buildNode(key string, m1, m2 map[string]interface{}) model.Node {
	v1, ok1 := m1[key]
	v2, ok2 := m2[key]

	if ok1 && ok2 {
		return buildNodeForBoth(key, v1, v2)
	}

	if ok1 {
		return model.Node{
			Key:   key,
			Type:  model.NodeRemoved,
			Value: v1,
		}
	}

	// !ok1 && ok2
	return model.Node{
		Key:   key,
		Type:  model.NodeAdded,
		Value: v2,
	}
}

func buildNodeForBoth(key string, v1, v2 interface{}) model.Node {
	m1Child, isMap1 := utils.ToMap(v1)
	m2Child, isMap2 := utils.ToMap(v2)

	if isMap1 && isMap2 {
		return model.Node{
			Key:      key,
			Type:     model.NodeNested,
			Children: buildAST(m1Child, m2Child),
		}
	}

	if reflect.DeepEqual(v1, v2) {
		return model.Node{
			Key:   key,
			Type:  model.NodeUnchanged,
			Value: v1,
		}
	}

	return model.Node{
		Key:      key,
		Type:     model.NodeUpdated,
		OldValue: v1,
		NewValue: v2,
	}
}
