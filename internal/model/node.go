package model

type NodeType string

const (
	NodeAdded     NodeType = "added"
	NodeRemoved   NodeType = "removed"
	NodeUnchanged NodeType = "unchanged"
	NodeUpdated   NodeType = "updated"
	NodeNested    NodeType = "nested"
)

type Node struct {
	Key      string
	Type     NodeType
	Value    interface{} // для added/removed/unchanged
	OldValue interface{} // для updated
	NewValue interface{} // для updated
	Children []Node      // для nested
}
