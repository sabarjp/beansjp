---
# beans-3
title: 'TUI: folded children are invisible to the text filter'
status: todo
type: bug
priority: low
created_at: 2026-08-28T04:47:04Z
updated_at: 2026-08-28T04:47:04Z
parent: beans-2
---

When a parent is folded (manually or automatically), its children are not in the Bubbletea list's item slice, so pressing `/` and typing cannot match them. Consider temporarily unfolding everything while a text filter is active.
