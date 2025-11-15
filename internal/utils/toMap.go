package utils

import "fmt"

func ToMap(v interface{}) (map[string]interface{}, bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		return convertStringMap(m), true
	case map[interface{}]interface{}:
		return convertInterfaceMap(m), true
	default:
		return nil, false
	}
}

func convertStringMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(m))
	for k, val := range m {
		res[k] = convertValue(val)
	}
	return res
}

func convertInterfaceMap(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(m))
	for k, val := range m {
		keyStr, ok := k.(string)
		if !ok {
			keyStr = fmt.Sprint(k)
		}
		res[keyStr] = convertValue(val)
	}
	return res
}

func convertValue(val interface{}) interface{} {
	if nested, ok := ToMap(val); ok {
		return nested
	}
	return val
}
