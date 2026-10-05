package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderBeanRow_NarrowWidth(t *testing.T) {
	// Test that RenderBeanRow doesn't panic with very small MaxTitleWidth values
	// This was a bug where MaxTitleWidth < 4 caused a slice bounds panic

	tests := []struct {
		name          string
		maxTitleWidth int
		title         string
	}{
		{"zero width", 0, "Test Title"},
		{"width 1", 1, "Test Title"},
		{"width 2", 2, "Test Title"},
		{"width 3", 3, "Test Title"},
		{"width 4", 4, "Test Title"},
		{"width 5", 5, "Test Title"},
		{"short title fits", 10, "Hi"},
		{"exact fit", 10, "0123456789"},
		{"needs truncation", 10, "This is a longer title"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("RenderBeanRow panicked with MaxTitleWidth=%d: %v", tt.maxTitleWidth, r)
				}
			}()

			cfg := BeanRowConfig{
				MaxTitleWidth: tt.maxTitleWidth,
				StatusColor:   "green",
				TypeColor:     "blue",
			}

			result := RenderBeanRow("abc123", "todo", "task", tt.title, cfg)
			if result == "" {
				t.Error("expected non-empty result")
			}
		})
	}
}

func TestRenderBeanRow_NarrowWidthWithPriority(t *testing.T) {
	// Priority symbol takes 2 extra chars, which reduces available title width
	// This tests that the adjustment doesn't cause negative slice bounds

	tests := []struct {
		name          string
		maxTitleWidth int
		priority      string
	}{
		{"width 1 with priority", 1, "high"},
		{"width 2 with priority", 2, "high"},
		{"width 3 with priority", 3, "critical"},
		{"width 4 with priority", 4, "high"},
		{"width 5 with priority", 5, "low"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("RenderBeanRow panicked with MaxTitleWidth=%d and priority=%s: %v",
						tt.maxTitleWidth, tt.priority, r)
				}
			}()

			cfg := BeanRowConfig{
				MaxTitleWidth: tt.maxTitleWidth,
				Priority:      tt.priority,
				PriorityColor: "red",
				StatusColor:   "green",
				TypeColor:     "blue",
			}

			result := RenderBeanRow("abc123", "todo", "task", "Long title that needs truncation", cfg)
			if result == "" {
				t.Error("expected non-empty result")
			}
		})
	}
}

func TestRenderBeanRow_ImplicitStatusFitsTitleColumn(t *testing.T) {
	// The "↑<status>" annotation must fit inside the title column budget,
	// otherwise long titles push the row past the terminal width and wrap.
	const longTitle = "A fairly long bean title that will certainly need truncating"

	// Fixed columns preceding the title: ID + type + status + separating spaces.
	const fixedWidth = ColWidthID + 1 + ColWidthType + 1 + ColWidthStatus + 1

	widths := []int{10, 15, 20, 24, 30, 40, 60}
	for _, w := range widths {
		for _, showTags := range []bool{false, true} {
			cfg := BeanRowConfig{
				MaxTitleWidth:  w,
				StatusColor:    "green",
				TypeColor:      "blue",
				Priority:       "high",
				PriorityColor:  "red",
				ImplicitStatus: "completed",
				ShowTags:       showTags,
				Tags:           []string{"idea"},
			}

			row := RenderBeanRow("abc123", "todo", "task", longTitle, cfg)

			want := fixedWidth + w
			if showTags {
				want += 1 + ColWidthTags
			}
			if got := lipgloss.Width(row); got > want {
				t.Errorf("MaxTitleWidth=%d showTags=%v: row width %d exceeds budget %d\nrow: %q",
					w, showTags, got, want, row)
			}
			// The annotation is only shown when it leaves a usable title behind.
			annotationFits := w-2-len([]rune(" ↑completed")) >= minTitleWidthForImplicit
			if got := strings.Contains(row, "↑completed"); got != annotationFits {
				t.Errorf("MaxTitleWidth=%d showTags=%v: annotation present=%v, want %v: %q",
					w, showTags, got, annotationFits, row)
			}
		}
	}
}

func TestRenderBeanRow_TruncatesByRunes(t *testing.T) {
	// Truncation must count runes, not bytes, so multi-byte titles aren't
	// cut short (or split mid-rune).
	title := "日本語のタイトルはとても長いです"

	cfg := BeanRowConfig{MaxTitleWidth: 10, StatusColor: "green", TypeColor: "blue"}
	row := RenderBeanRow("abc123", "todo", "task", title, cfg)

	want := string([]rune(title)[:7]) + "..."
	if !strings.Contains(row, want) {
		t.Errorf("expected row to contain %q, got %q", want, row)
	}
	if !utf8ValidRow(row) {
		t.Errorf("row contains invalid UTF-8: %q", row)
	}
}

func utf8ValidRow(s string) bool {
	return strings.ToValidUTF8(s, "�") == s
}

func TestShortType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"milestone", "M"},
		{"epic", "E"},
		{"bug", "B"},
		{"feature", "F"},
		{"task", "T"},
		{"unknown", "?"},
		{"", "?"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ShortType(tt.input)
			if result != tt.expected {
				t.Errorf("ShortType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestShortStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"draft", "D"},
		{"todo", "T"},
		{"in-progress", "I"},
		{"stale", "X"},
		{"completed", "C"},
		{"scrapped", "S"},
		{"unknown", "?"},
		{"", "?"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ShortStatus(tt.input)
			if result != tt.expected {
				t.Errorf("ShortStatus(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
