package parser

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func fixturePath(name string) string {
	return filepath.Join("..", "..", "testdata", "fixture", name)
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"json", "config.json", "json"},
		{"yaml", "config.yaml", "yaml"},
		{"yml", "config.yml", "yaml"},
		{"upper_json", "CONFIG.JSON", "json"},
		{"unknown", "config.txt", ""},
		{"no_ext", "config", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := detectFormat(tt.path)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestParse_JSONSuccess(t *testing.T) {
	path := fixturePath("valid.json")

	result, err := Parse(path)
	require.NoError(t, err)

	require.Equal(t, "value", result["key"])
	require.Equal(t, float64(42), result["num"])
}

func TestParse_YAMLSuccess(t *testing.T) {
	path := fixturePath("valid.yaml")

	result, err := Parse(path)
	require.NoError(t, err)

	require.Equal(t, "value", result["key"])
	require.Equal(t, "42", fmt.Sprint(result["num"]))

	nested, ok := result["nested"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, true, nested["flag"])
}

func TestParse_UnsupportedFormat(t *testing.T) {
	path := fixturePath("unsupported.txt")

	_, err := Parse(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported file format")
}

func TestParse_FileDoesNotExist(t *testing.T) {
	path := fixturePath("does-not-exist.json")

	_, err := Parse(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot read file")
}

func TestParse_InvalidJSON(t *testing.T) {
	path := fixturePath("invalid.json")

	_, err := Parse(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid json")
}

func TestParse_InvalidYAML(t *testing.T) {
	path := fixturePath("invalid.yaml")

	_, err := Parse(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid yaml")
}
