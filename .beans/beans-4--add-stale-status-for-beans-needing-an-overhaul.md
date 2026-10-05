---
# beans-4
title: Add 'stale' status for beans needing an overhaul
status: completed
type: feature
priority: normal
created_at: 2026-10-05T18:38:23Z
updated_at: 2026-10-05T18:47:15Z
---

Add a non-terminal `stale` status that sorts between draft and completed, for beans whose existing content needs an overhaul.

- [x] Add stale to DefaultStatuses (between draft and completed)
- [x] Exclude stale from --ready
- [x] TUI short code + styling
- [x] Web UI: status order, form dropdown, colors, workflow actions
- [x] Docs/schema descriptions
- [x] Tests (config, list --ready, e2e)

## Summary of Changes

- New non-terminal `stale` status (orange), ordered in-progress → todo → draft → stale → completed → scrapped.
- `beans list --ready` excludes stale; default `beans list` still shows it. Filter logic extracted to `applyReadyFilter` with a test.
- Not archivable and still counts as an active blocker (same as draft).
- TUI: short code `X` (S is scrapped), new `orange` named color.
- Web UI: Stale section in the backlog under Draft; status dropdown, badge colors (light/dark), frontend sort order.
- Workflow actions: draft → Todo/Stale/Scrap, todo → Stale/Scrap, stale → Draft/Scrap.
- Prime template, GraphQL schema description, and generated types updated.
- Tests: config/ui/links unit tests, ready-filter test, e2e for backlog ordering and stale workflow buttons.
