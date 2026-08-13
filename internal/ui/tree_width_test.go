package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/hmans/beans/pkg/bean"
	"github.com/hmans/beans/pkg/config"
)

// TestRenderBeanRow_FitsResponsiveColumnBudget mirrors the width math the TUI
// list delegate uses (internal/tui/list.go) and asserts a rendered row never
// exceeds the viewport width — including at widths where full type/status
// names and the tags column kick in.
func TestRenderBeanRow_FitsResponsiveColumnBudget(t *testing.T) {
	for width := 60; width <= 220; width++ {
		cols := CalculateResponsiveColumns(width, true)

		baseWidth := cols.ID + cols.Status + cols.Type + 4 // 4 for cursor + padding
		if cols.ShowTags {
			baseWidth += cols.Tags + 1
		}
		maxTitleWidth := max(0, width-baseWidth)

		row := RenderBeanRow("bean-abcd", "in-progress", "milestone",
			"A long bean title that will fill the entire available title column", BeanRowConfig{
				StatusColor:    "green",
				TypeColor:      "blue",
				Priority:       "high",
				PriorityColor:  "red",
				MaxTitleWidth:  maxTitleWidth,
				ShowCursor:     true,
				Tags:           []string{"idea", "ui"},
				ShowTags:       cols.ShowTags,
				TagsColWidth:   cols.Tags,
				MaxTags:        cols.MaxTags,
				UseFullNames:   cols.UseFullTypeStatus,
				ImplicitStatus: "completed",
			})

		if got := lipgloss.Width(row); got > width {
			t.Fatalf("viewport width %d: row width %d exceeds it: %q", width, got, row)
		}
	}
}

// TestRenderTree_NoLineExceedsTermWidth guards against rendered rows spilling
// past the terminal width, which makes the terminal wrap them onto a second
// line and garbles the list.
func TestRenderTree_NoLineExceedsTermWidth(t *testing.T) {
	parent := &bean.Bean{
		ID:     "bean-parent",
		Title:  "A parent bean with a reasonably long title for testing purposes",
		Status: "completed",
		Type:   "epic",
		Tags:   []string{"idea"},
	}
	child := &bean.Bean{
		ID:       "bean-child",
		Title:    "A child bean whose title is long enough to fill the whole title column",
		Status:   "todo",
		Type:     "task",
		Priority: "high",
		Tags:     []string{"idea", "ui"},
	}

	nodes := []*TreeNode{{
		Bean:    parent,
		Matched: true,
		Children: []*TreeNode{{
			Bean:           child,
			Matched:        true,
			ImplicitStatus: "completed",
		}},
	}}

	cfg := config.Default()

	for termWidth := 60; termWidth <= 200; termWidth++ {
		out := RenderTree(nodes, cfg, 11, true, termWidth)
		for i, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			if got := lipgloss.Width(line); got > termWidth {
				t.Fatalf("termWidth=%d line %d width %d exceeds terminal width: %q",
					termWidth, i, got, line)
			}
		}
	}
}
