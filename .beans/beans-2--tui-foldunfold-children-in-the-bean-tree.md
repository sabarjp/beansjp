---
# beans-2
title: 'TUI: fold/unfold children in the bean tree'
status: completed
type: feature
priority: normal
created_at: 2026-08-28T04:38:49Z
updated_at: 2026-08-28T04:47:27Z
---

Add directory-tree style folding to the TUI list. Toggle a parent's children with a key; parents whose descendants are all complete/scrapped start auto-folded.

## Todo

- [x] Add collapse support to ui.FlattenTree (hidden counts, has-children flags)
- [x] Auto-fold parents whose subtrees are entirely terminal
- [x] TUI fold state + key bindings (tab / left / right)
- [x] Render fold markers in the list rows
- [x] Tests

## Summary of Changes

- `internal/ui/tree.go`: `FlattenTreeCollapsed` omits the children of folded nodes and reports `HasChildren`/`Collapsed`/`HiddenCount` per row; `AutoCollapsed` returns the parents whose children are all in a terminal status (recursively, honouring implicit status); `FoldMarker` renders the fixed-width open/closed triangle column.
- `internal/tui/list.go`: fold state on the list model (automatic state plus explicit user overrides that win), `tab` to toggle, left/`h` to fold or jump to the parent, right/`l` to unfold or step into the first child. Left/right no longer page the list. The ID column reserves the marker width only when something is foldable, and folded rows show `(+n)` hidden descendants.
- Help overlay and footer mention the new keys.
- Tests: `internal/ui/tree_fold_test.go`, `internal/tui/fold_test.go`.

Follow-up: [[beans-3]] - folded children cannot be matched by the text filter.
