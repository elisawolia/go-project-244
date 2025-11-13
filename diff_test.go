package code

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const expectedStylish = `{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow: 
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: {
            key: value
        }
      + nest: str
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}`

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

func TestGenDiff_YAMLFlat(t *testing.T) {
	path1 := fixturePath("file1.yaml")
	path2 := fixturePath("file2.yaml")

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

func TestGenDiffStylishJSONested(t *testing.T) {
	path1 := fixturePath("file_hard1.json")
	path2 := fixturePath("file_hard2.json")

	got, err := GenDiff(path1, path2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != expectedStylish {
		t.Errorf("result mismatch.\nExpected:\n%v\n\nGot:\n%v", expectedStylish, got)
	}
}

func TestGenDiff_FileNotFound(t *testing.T) {
	path1 := fixturePath("simple1.json")
	path2 := fixturePath("does-not-exist.json")

	_, err := GenDiff(path1, path2)
	require.Error(t, err)
}
