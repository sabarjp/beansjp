---
# beans-3
title: 'TUI: folded children are invisible to the text filter'
status: completed
type: bug
priority: low
created_at: 2026-08-28T04:47:04Z
updated_at: 2026-08-28T04:52:16Z
parent: beans-2
---

When a parent is folded (manually or automatically), its children are not in the Bubbletea list's item slice, so pressing `/` and typing cannot match them. Consider temporarily unfolding everything while a text filter is active.

## Summary of Changes

Folds are now suspended while a text filter is active, so every bean is in the list and can be matched.

- `listModel.foldsSuspended` short-circuits `effectiveFolds`; entering or leaving the filter (detected as a transition in/out of `list.Unfiltered`) triggers a reload, and the fold keys are inert while suspended.
- The command returned by `list.SetItems` is now forwarded, so an applied filter re-runs over the rebuilt items instead of going stale.
- Cursor restoration remembers the selected bean *and* its ancestors, taken from before the update (clearing a filter moves the cursor first), and falls back to the nearest visible ancestor when the bean itself is folded away.
- `selectBeanID` indexes `VisibleItems` rather than `Items`, which is what list indices refer to under a filter.
- The preview pane now follows the selected bean ID instead of the cursor index, so it keeps up when filtering changes the selection without moving the cursor.

Verified in a real TUI session: filtering matches a bean hidden under an auto-folded parent, clearing the filter restores the folds and parks the cursor on the folded ancestor, and manual unfolds survive the round-trip.
