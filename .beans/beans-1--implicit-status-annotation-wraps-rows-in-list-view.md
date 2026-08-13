---
# beans-1
title: Implicit status annotation wraps rows in list views
status: completed
type: bug
created_at: 2026-08-13T01:21:39Z
updated_at: 2026-08-13T01:21:39Z
---

The " ↑<status>" implicit-status annotation was appended after the title/tags columns without being counted in the width budget, so at certain terminal widths rows exceeded the terminal width and wrapped.

## Summary of Changes

- `RenderBeanRow` now reserves the annotation's width inside the title column (title truncates to make room), and drops the annotation entirely when fewer than 4 title chars would remain.
- Annotation moved before the tags column so tags stay aligned.
- Title truncation is now rune-based instead of byte-based (multi-byte titles were truncated early / could split a rune).
- Fixed an off-by-one in the tags-column width budget (`RenderTree` and the TUI list delegate did not account for the space separating title and tags), which made every row one column too wide at 140+ columns.
- Full type/status column widths are now shared constants (`ColWidthFullType`/`ColWidthFullStatus`); `CalculateResponsiveColumns` reported 10 for type while the renderer used 12, overflowing the TUI list at 120+ columns.
- Added tests asserting no rendered row exceeds the terminal width across widths 60-220.
