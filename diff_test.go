package code

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func fixturePath(name string) string {
	return filepath.Join("testdata", "fixture", name)
}

func TestGenDiff_JSONFlat(t *testing.T) {
	path1 := fixturePath("file1.json")
	path2 := fixturePath("file2.json")

	actual, err := GenDiff(path1, path2)
	require.NoError(t, err)

	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`

	require.Equal(t, expected, actual)
}

func TestGenDiff_FileNotFound(t *testing.T) {
	path1 := fixturePath("simple1.json")
	path2 := fixturePath("does-not-exist.json")

	_, err := GenDiff(path1, path2)
	require.Error(t, err)
}
