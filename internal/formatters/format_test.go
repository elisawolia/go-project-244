package formatters

import (
	"testing"

	"code/internal/model"
)

func TestFormat_DefaultsToStylish(t *testing.T) {
	var tree []model.Node

	want := formatStylish(tree)

	got, err := Format(tree, "")
	if err != nil {
		t.Fatalf("Format returned error for empty format: %v", err)
	}

	if got != want {
		t.Fatalf("Format() with empty format returned %q, want %q", got, want)
	}
}

func TestFormat_Stylish(t *testing.T) {
	var tree []model.Node

	want := formatStylish(tree)

	got, err := Format(tree, "stylish")
	if err != nil {
		t.Fatalf("Format returned error for stylish: %v", err)
	}

	if got != want {
		t.Fatalf("Format(\"stylish\") = %q, want %q", got, want)
	}
}

func TestFormat_Plain(t *testing.T) {
	var tree []model.Node

	want := formatPlain(tree)

	got, err := Format(tree, "plain")
	if err != nil {
		t.Fatalf("Format returned error for plain: %v", err)
	}

	if got != want {
		t.Fatalf("Format(\"plain\") = %q, want %q", got, want)
	}
}

func TestFormat_JSON(t *testing.T) {
	var tree []model.Node

	want, wantErr := formatJSON(tree)

	got, err := Format(tree, "json")

	if (err != nil) != (wantErr != nil) {
		t.Fatalf("Format(\"json\") error = %v, want %v", err, wantErr)
	}

	if err != nil && err.Error() != wantErr.Error() {
		t.Fatalf("Format(\"json\") error = %v, want %v", err, wantErr)
	}

	if got != want {
		t.Fatalf("Format(\"json\") = %q, want %q", got, want)
	}
}

func TestFormat_UnknownFormat(t *testing.T) {
	var tree []model.Node

	_, err := Format(tree, "xml")
	if err == nil {
		t.Fatalf("expected error for unknown format, got nil")
	}

	wantMsg := "unknown format: xml"
	if err.Error() != wantMsg {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), wantMsg)
	}
}

func TestFormatJSON_EmptyTree(t *testing.T) {
	t.Parallel()

	tree := []model.Node{}

	got, err := formatJSON(tree)
	if err != nil {
		t.Fatalf("formatJSON returned error for empty tree: %v", err)
	}

	want := "[]"
	if got != want {
		t.Fatalf("formatJSON = %q, want %q", got, want)
	}
}

func TestFormatJSON_Tree(t *testing.T) {
	t.Parallel()

	tree := []model.Node{
		{
			Key:   "key",
			Type:  model.NodeAdded,
			Value: "value",
		},
		{
			Key:      "changed_key",
			Type:     model.NodeUpdated,
			OldValue: "old_value",
			NewValue: "new_value",
		},
	}

	got, err := formatJSON(tree)
	if err != nil {
		t.Fatalf("formatJSON returned error for empty tree: %v", err)
	}

	want := `[
  {
    "Key": "key",
    "Type": "added",
    "Value": "value",
    "OldValue": null,
    "NewValue": null,
    "Children": null
  },
  {
    "Key": "changed_key",
    "Type": "updated",
    "Value": null,
    "OldValue": "old_value",
    "NewValue": "new_value",
    "Children": null
  }
]`
	if got != want {
		t.Fatalf("formatJSON = %q, want %q", got, want)
	}
}

func TestFormatPlain_EmptyTree(t *testing.T) {
	t.Parallel()

	tree := []model.Node{}

	got := formatPlain(tree)

	want := ""
	if got != want {
		t.Fatalf("formatPlain = %q, want %q", got, want)
	}
}

func TestFormatPlain_Tree(t *testing.T) {
	t.Parallel()

	tree := []model.Node{
		{
			Key:  "group",
			Type: model.NodeNested,
			Children: []model.Node{
				{
					Key:   "added",
					Type:  model.NodeAdded,
					Value: "val",
				},
				{
					Key:  "removed",
					Type: model.NodeRemoved,
				},
				{
					Key:      "updated",
					Type:     model.NodeUpdated,
					OldValue: "old",
					NewValue: "new",
				},
				{
					Key:   "complexMap",
					Type:  model.NodeAdded,
					Value: map[string]interface{}{"a": 1},
				},
				{
					Key:   "complexList",
					Type:  model.NodeAdded,
					Value: []interface{}{1, 2},
				},
				{
					Key:   "unchangedInside",
					Type:  model.NodeUnchanged,
					Value: "noop",
				},
			},
		},
		{
			Key:   "flag",
			Type:  model.NodeAdded,
			Value: true,
		},
		{
			Key:   "disabled",
			Type:  model.NodeAdded,
			Value: false,
		},
		{
			Key:      "count",
			Type:     model.NodeUpdated,
			OldValue: 10,
			NewValue: nil,
		},
		{
			Key:   "unchanged",
			Type:  model.NodeUnchanged,
			Value: "noop",
		},
	}

	got := formatPlain(tree)

	want := `Property 'group.added' was added with value: 'val'
Property 'group.removed' was removed
Property 'group.updated' was updated. From 'old' to 'new'
Property 'group.complexMap' was added with value: [complex value]
Property 'group.complexList' was added with value: [complex value]
Property 'flag' was added with value: true
Property 'disabled' was added with value: false
Property 'count' was updated. From 10 to null`

	if got != want {
		t.Fatalf("formatPlain = %q, want %q", got, want)
	}
}

func TestFormatStylish_Tree(t *testing.T) {
	t.Parallel()

	tree := []model.Node{
		{
			Key:  "group",
			Type: model.NodeNested,
			Children: []model.Node{
				{
					Key:   "added",
					Type:  model.NodeAdded,
					Value: "val",
				},
				{
					Key:  "removed",
					Type: model.NodeRemoved,
				},
				{
					Key:      "updated",
					Type:     model.NodeUpdated,
					OldValue: "old",
					NewValue: "new",
				},
				{
					Key:   "complexMap",
					Type:  model.NodeAdded,
					Value: map[string]interface{}{"a": 1},
				},
				{
					Key:   "complexList",
					Type:  model.NodeAdded,
					Value: []interface{}{1, 2},
				},
				{
					Key:   "unchangedInside",
					Type:  model.NodeUnchanged,
					Value: "noop",
				},
			},
		},
		{
			Key:   "flag",
			Type:  model.NodeAdded,
			Value: true,
		},
		{
			Key:   "disabled",
			Type:  model.NodeAdded,
			Value: false,
		},
		{
			Key:      "count",
			Type:     model.NodeUpdated,
			OldValue: 10,
			NewValue: nil,
		},
		{
			Key:   "unchanged",
			Type:  model.NodeUnchanged,
			Value: "noop",
		},
	}

	got := formatStylish(tree)

	want := `{
    group: {
      + added: val
      - removed: null
      - updated: old
      + updated: new
      + complexMap: {
            a: 1
        }
      + complexList: [
            1
            2
        ]
        unchangedInside: noop
    }
  + flag: true
  + disabled: false
  - count: 10
  + count: null
    unchanged: noop
}`

	if got != want {
		t.Fatalf("formatStylish = %q, want %q", got, want)
	}
}
