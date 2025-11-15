package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"code/internal/utils"
)

func Parse(path string) (map[string]interface{}, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve path %q: %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %q: %w", absPath, err)
	}

	format := detectFormat(absPath)
	switch format {
	case "json":
		return parseJSON(data)
	case "yaml":
		return parseYAML(data)
	default:
		return nil, fmt.Errorf("unsupported file format for %q", absPath)
	}
}

func detectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return "json"
	case ".yml", ".yaml":
		return "yaml"
	default:
		return ""
	}
}

func parseJSON(b []byte) (map[string]interface{}, error) {
	var raw interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}
	return normalizeRoot(raw)
}

func parseYAML(b []byte) (map[string]interface{}, error) {
	var raw interface{}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("invalid yaml: %w", err)
	}
	return normalizeRoot(raw)
}

func normalizeRoot(raw interface{}) (map[string]interface{}, error) {
	if m, ok := utils.ToComplexStruct(raw); ok {
		if res, ok2 := m.(map[string]interface{}); ok2 {
			return res, nil
		}
	}
	if slice, ok := raw.([]interface{}); ok {
		res := make(map[string]interface{}, len(slice))
		for i, v := range slice {
			key := fmt.Sprintf("%d", i)
			res[key] = v
		}
		return res, nil
	}

	return nil, fmt.Errorf("unsupported root type %T", raw)
}
