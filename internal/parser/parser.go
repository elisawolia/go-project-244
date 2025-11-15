package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	if m, ok := utils.ToMap(raw); ok {
		return m, nil
	}

	if s, ok := raw.([]interface{}); ok {
		res := make(map[string]interface{}, len(s))
		for i, v := range s {
			res[strconv.Itoa(i)] = v
		}
		return res, nil
	}

	return map[string]interface{}{
		"__root__": raw,
	}, nil
}
