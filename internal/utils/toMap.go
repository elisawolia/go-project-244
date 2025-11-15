package utils

import "fmt"

func ToMap(v interface{}) (map[string]interface{}, bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{}, len(m))
		for k, val := range m {
			if nested, ok := ToMap(val); ok {
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
			if nested, ok := ToMap(val); ok {
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
