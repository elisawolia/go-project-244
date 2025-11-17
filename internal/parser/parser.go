package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}
	return result, nil
}

func parseYAML(b []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("invalid yaml: %w", err)
	}
	return result, nil
}
