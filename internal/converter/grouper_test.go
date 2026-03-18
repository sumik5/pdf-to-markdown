package converter

import (
	"testing"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

func TestGroupIntoLines(t *testing.T) {
	tests := []struct {
		name      string
		items     []model.TextItem
		wantLines int
		wantTexts []string
	}{
		{
			name:      "empty input",
			items:     nil,
			wantLines: 0,
			wantTexts: nil,
		},
		{
			name: "single item",
			items: []model.TextItem{
				{Text: "Hello", X: 10, Y: 100, FontSize: 12},
			},
			wantLines: 1,
			wantTexts: []string{"Hello"},
		},
		{
			name: "two items same Y",
			items: []model.TextItem{
				{Text: "World", X: 60, Y: 100, Width: 30, FontSize: 12},
				{Text: "Hello", X: 10, Y: 100, Width: 40, FontSize: 12},
			},
			wantLines: 1,
			wantTexts: []string{"Hello World"},
		},
		{
			name: "two items within Y tolerance",
			items: []model.TextItem{
				{Text: "Line1", X: 10, Y: 100.0, FontSize: 12},
				{Text: "same", X: 60, Y: 100.5, FontSize: 12},
			},
			wantLines: 1,
			wantTexts: []string{"Line1 same"},
		},
		{
			name: "two distinct lines",
			items: []model.TextItem{
				{Text: "Second", X: 10, Y: 80, FontSize: 12},
				{Text: "First", X: 10, Y: 100, FontSize: 12},
			},
			wantLines: 2,
			wantTexts: []string{"First", "Second"},
		},
		{
			name: "top to bottom ordering",
			items: []model.TextItem{
				{Text: "C", X: 10, Y: 60, FontSize: 12},
				{Text: "A", X: 10, Y: 100, FontSize: 12},
				{Text: "B", X: 10, Y: 80, FontSize: 12},
			},
			wantLines: 3,
			wantTexts: []string{"A", "B", "C"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupIntoLines(tt.items)
			if len(got) != tt.wantLines {
				t.Errorf("groupIntoLines() = %d lines, want %d", len(got), tt.wantLines)
				return
			}
			for i, wantText := range tt.wantTexts {
				if got[i].Text != wantText {
					t.Errorf("line[%d].Text = %q, want %q", i, got[i].Text, wantText)
				}
			}
		})
	}
}

func TestDetectIndentation(t *testing.T) {
	tests := []struct {
		name        string
		lines       []model.Line
		wantIndents []int
	}{
		{
			name:        "empty lines",
			lines:       nil,
			wantIndents: []int{},
		},
		{
			name: "no indentation",
			lines: []model.Line{
				{Items: []model.TextItem{{X: 10}}, Text: "A"},
				{Items: []model.TextItem{{X: 10}}, Text: "B"},
			},
			wantIndents: []int{0, 0},
		},
		{
			name: "one level indent",
			lines: []model.Line{
				{Items: []model.TextItem{{X: 10}}, Text: "A"},
				{Items: []model.TextItem{{X: 20}}, Text: "B"},
			},
			wantIndents: []int{0, 1},
		},
		{
			name: "two levels indent",
			lines: []model.Line{
				{Items: []model.TextItem{{X: 10}}, Text: "A"},
				{Items: []model.TextItem{{X: 30}}, Text: "B"},
			},
			wantIndents: []int{0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectIndentation(tt.lines)
			if len(got) != len(tt.wantIndents) {
				t.Fatalf("detectIndentation() = %d lines, want %d", len(got), len(tt.wantIndents))
			}
			for i, wantIndent := range tt.wantIndents {
				if got[i].IndentLevel != wantIndent {
					t.Errorf("line[%d].IndentLevel = %d, want %d", i, got[i].IndentLevel, wantIndent)
				}
			}
		})
	}
}

func TestJoinItemTexts(t *testing.T) {
	tests := []struct {
		name  string
		items []model.TextItem
		want  string
	}{
		{
			name:  "empty",
			items: nil,
			want:  "",
		},
		{
			name:  "single item",
			items: []model.TextItem{{Text: "Hello"}},
			want:  "Hello",
		},
		{
			name: "adjacent items with gap",
			items: []model.TextItem{
				{Text: "Hello", X: 10, Width: 30},
				{Text: "World", X: 45},
			},
			want: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinItemTexts(tt.items)
			if got != tt.want {
				t.Errorf("joinItemTexts() = %q, want %q", got, tt.want)
			}
		})
	}
}
