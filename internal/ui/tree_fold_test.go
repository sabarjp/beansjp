package ui

import (
	"testing"

	"github.com/hmans/beans/pkg/bean"
)

// foldTestTree builds:
// m1 (todo)
//
//	├── e1 (todo)
//	│    └── t1 (todo)
//	└── e2 (completed)
//	     └── t2 (completed)
func foldTestTree(t *testing.T) []*TreeNode {
	t.Helper()

	m1 := &bean.Bean{ID: "m1", Title: "Milestone 1", Type: "milestone", Status: "todo"}
	e1 := &bean.Bean{ID: "e1", Title: "Epic 1", Type: "epic", Status: "todo", Parent: "m1"}
	t1 := &bean.Bean{ID: "t1", Title: "Task 1", Type: "task", Status: "todo", Parent: "e1"}
	e2 := &bean.Bean{ID: "e2", Title: "Epic 2", Type: "epic", Status: "completed", Parent: "m1"}
	t2 := &bean.Bean{ID: "t2", Title: "Task 2", Type: "task", Status: "completed", Parent: "e2"}

	all := []*bean.Bean{m1, e1, t1, e2, t2}
	noSort := func(b []*bean.Bean) {}
	return BuildTree(all, all, noSort, nil)
}

func flatIDs(items []FlatItem) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.Bean.ID
	}
	return ids
}

func TestFlattenTreeCollapsed(t *testing.T) {
	tree := foldTestTree(t)

	t.Run("nil collapse set shows everything", func(t *testing.T) {
		items := FlattenTreeCollapsed(tree, nil)
		if len(items) != 5 {
			t.Fatalf("expected 5 items, got %d (%v)", len(items), flatIDs(items))
		}
		for _, item := range items {
			if item.Collapsed {
				t.Errorf("%s: expected no item to be collapsed", item.Bean.ID)
			}
		}
	})

	t.Run("collapsed node hides its descendants", func(t *testing.T) {
		items := FlattenTreeCollapsed(tree, map[string]bool{"e1": true})

		ids := flatIDs(items)
		want := []string{"m1", "e1", "e2", "t2"}
		if len(ids) != len(want) {
			t.Fatalf("got %v, want %v", ids, want)
		}
		for i, id := range want {
			if ids[i] != id {
				t.Fatalf("got %v, want %v", ids, want)
			}
		}

		for _, item := range items {
			if item.Bean.ID != "e1" {
				continue
			}
			if !item.Collapsed {
				t.Error("e1: expected Collapsed = true")
			}
			if item.HiddenCount != 1 {
				t.Errorf("e1: HiddenCount = %d, want 1", item.HiddenCount)
			}
		}
	})

	t.Run("collapsing a root hides the whole subtree", func(t *testing.T) {
		items := FlattenTreeCollapsed(tree, map[string]bool{"m1": true})
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %v", flatIDs(items))
		}
		if items[0].HiddenCount != 4 {
			t.Errorf("m1: HiddenCount = %d, want 4", items[0].HiddenCount)
		}
	})

	t.Run("marks nodes that have children", func(t *testing.T) {
		items := FlattenTreeCollapsed(tree, nil)
		wantChildren := map[string]bool{"m1": true, "e1": true, "e2": true}
		for _, item := range items {
			if got := item.HasChildren; got != wantChildren[item.Bean.ID] {
				t.Errorf("%s: HasChildren = %v, want %v", item.Bean.ID, got, wantChildren[item.Bean.ID])
			}
		}
	})

	t.Run("collapsed flag is ignored for leaves", func(t *testing.T) {
		items := FlattenTreeCollapsed(tree, map[string]bool{"t1": true})
		if len(items) != 5 {
			t.Fatalf("expected 5 items, got %v", flatIDs(items))
		}
		for _, item := range items {
			if item.Bean.ID == "t1" && item.Collapsed {
				t.Error("t1: leaf should never report as collapsed")
			}
		}
	})
}

func TestAutoCollapsed(t *testing.T) {
	isDone := func(status string) bool {
		return status == "completed" || status == "scrapped"
	}

	t.Run("folds parents whose subtree is entirely done", func(t *testing.T) {
		folded := AutoCollapsed(foldTestTree(t), isDone)

		if !folded["e2"] {
			t.Error("e2: expected auto-fold (all descendants completed)")
		}
		if folded["e1"] {
			t.Error("e1: expected no auto-fold (open child)")
		}
		if folded["m1"] {
			t.Error("m1: expected no auto-fold (open descendant)")
		}
	})

	t.Run("does not fold leaves", func(t *testing.T) {
		b := &bean.Bean{ID: "t1", Title: "Task 1", Status: "completed"}
		all := []*bean.Bean{b}
		tree := BuildTree(all, all, func([]*bean.Bean) {}, nil)

		if len(AutoCollapsed(tree, isDone)) != 0 {
			t.Error("expected no auto-folds for a childless bean")
		}
	})

	t.Run("folds when the parent itself is still open", func(t *testing.T) {
		parent := &bean.Bean{ID: "e1", Title: "Epic 1", Status: "in-progress"}
		child := &bean.Bean{ID: "t1", Title: "Task 1", Status: "completed", Parent: "e1"}
		all := []*bean.Bean{parent, child}
		tree := BuildTree(all, all, func([]*bean.Bean) {}, nil)

		if !AutoCollapsed(tree, isDone)["e1"] {
			t.Error("e1: expected auto-fold")
		}
	})

	t.Run("an open parent of an open child does not fold", func(t *testing.T) {
		parent := &bean.Bean{ID: "e1", Title: "Epic 1", Status: "in-progress"}
		done := &bean.Bean{ID: "t1", Title: "Task 1", Status: "completed", Parent: "e1"}
		open := &bean.Bean{ID: "t2", Title: "Task 2", Status: "todo", Parent: "e1"}
		all := []*bean.Bean{parent, done, open}
		tree := BuildTree(all, all, func([]*bean.Bean) {}, nil)

		if AutoCollapsed(tree, isDone)["e1"] {
			t.Error("e1: expected no auto-fold")
		}
	})

	t.Run("uses implicit status inherited from ancestors", func(t *testing.T) {
		parent := &bean.Bean{ID: "e1", Title: "Epic 1", Status: "completed"}
		child := &bean.Bean{ID: "t1", Title: "Task 1", Status: "todo", Parent: "e1"}
		all := []*bean.Bean{parent, child}
		implicit := map[string]string{"t1": "completed"}
		tree := BuildTree(all, all, func([]*bean.Bean) {}, implicit)

		if !AutoCollapsed(tree, isDone)["e1"] {
			t.Error("e1: expected auto-fold (child inherits completed)")
		}
	})
}

func TestFoldMarker(t *testing.T) {
	tests := []struct {
		name        string
		hasChildren bool
		collapsed   bool
		want        string
	}{
		{"leaf", false, false, foldMarkerNone},
		{"open parent", true, false, foldMarkerOpen},
		{"folded parent", true, true, foldMarkerClosed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FoldMarker(tt.hasChildren, tt.collapsed); got != tt.want {
				t.Errorf("FoldMarker(%v, %v) = %q, want %q", tt.hasChildren, tt.collapsed, got, tt.want)
			}
			if len([]rune(FoldMarker(tt.hasChildren, tt.collapsed))) != FoldMarkerWidth {
				t.Errorf("marker width != FoldMarkerWidth")
			}
		})
	}
}
