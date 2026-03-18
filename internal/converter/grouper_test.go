package converter

import (
	"testing"
)

func TestParseLinesWithIndentation(t *testing.T) {
	tests := []struct {
		name        string
		rawLines    []string
		wantTexts   []string
		wantIndents []int
	}{
		{
			name:        "empty input",
			rawLines:    nil,
			wantTexts:   nil,
			wantIndents: nil,
		},
		{
			name:        "single line",
			rawLines:    []string{"Hello World"},
			wantTexts:   []string{"Hello World"},
			wantIndents: []int{0},
		},
		{
			name:        "blank lines are skipped",
			rawLines:    []string{"First", "", "  ", "\t", "Second"},
			wantTexts:   []string{"First", "Second"},
			wantIndents: []int{0, 0},
		},
		{
			name:        "trailing whitespace is removed",
			rawLines:    []string{"Hello   "},
			wantTexts:   []string{"Hello"},
			wantIndents: []int{0},
		},
		{
			name:        "no indentation",
			rawLines:    []string{"Line A", "Line B"},
			wantTexts:   []string{"Line A", "Line B"},
			wantIndents: []int{0, 0},
		},
		{
			name:        "one level indent (4 spaces)",
			rawLines:    []string{"Base", "    Indented"},
			wantTexts:   []string{"Base", "Indented"},
			wantIndents: []int{0, 1},
		},
		{
			name:        "two levels indent (8 spaces)",
			rawLines:    []string{"Base", "        Deep"},
			wantTexts:   []string{"Base", "Deep"},
			wantIndents: []int{0, 2},
		},
		{
			name:        "all lines equally indented normalizes to 0",
			rawLines:    []string{"    Line A", "    Line B"},
			wantTexts:   []string{"Line A", "Line B"},
			wantIndents: []int{0, 0},
		},
		{
			name:        "relative indent from minimum",
			rawLines:    []string{"    Base", "        Indented"},
			wantTexts:   []string{"Base", "Indented"},
			wantIndents: []int{0, 1},
		},
		{
			name:        "mixed indent levels",
			rawLines:    []string{"Top", "    Child", "        Grandchild", "    Child2"},
			wantTexts:   []string{"Top", "Child", "Grandchild", "Child2"},
			wantIndents: []int{0, 1, 2, 1},
		},
		{
			name:        "normalizes multiple spaces",
			rawLines:    []string{"Programming                          Language   Pragmatics"},
			wantTexts:   []string{"Programming Language Pragmatics"},
			wantIndents: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLinesWithIndentation(tt.rawLines)

			if tt.wantTexts == nil {
				if got != nil {
					t.Errorf("parseLinesWithIndentation() = %v, want nil", got)
				}
				return
			}

			if len(got) != len(tt.wantTexts) {
				t.Fatalf("parseLinesWithIndentation() = %d lines, want %d", len(got), len(tt.wantTexts))
			}

			for i, wantText := range tt.wantTexts {
				if got[i].Text != wantText {
					t.Errorf("line[%d].Text = %q, want %q", i, got[i].Text, wantText)
				}
				if got[i].IndentLevel != tt.wantIndents[i] {
					t.Errorf("line[%d].IndentLevel = %d, want %d", i, got[i].IndentLevel, tt.wantIndents[i])
				}
			}
		})
	}
}
