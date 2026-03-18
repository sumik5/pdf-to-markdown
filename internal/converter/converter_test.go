package converter

import (
	"testing"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

func TestProcessPage(t *testing.T) {
	tests := []struct {
		name       string
		lines      []string
		pageNumber int
		outline    []model.OutlineItem
		wantLines  int
	}{
		{
			name:       "empty page",
			lines:      nil,
			pageNumber: 1,
			outline:    nil,
			wantLines:  0,
		},
		{
			name:       "single line page",
			lines:      []string{"Hello World"},
			pageNumber: 1,
			outline:    nil,
			wantLines:  1,
		},
		{
			name:       "blank lines are skipped",
			lines:      []string{"First line", "", "  ", "Second line"},
			pageNumber: 2,
			outline:    nil,
			wantLines:  2,
		},
		{
			name:       "heading detected via outline",
			lines:      []string{"Introduction", "Normal text here."},
			pageNumber: 1,
			outline:    []model.OutlineItem{{Title: "Introduction"}},
			wantLines:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessPage(tt.lines, tt.pageNumber, tt.outline)
			if got.Number != tt.pageNumber {
				t.Errorf("Page.Number = %d, want %d", got.Number, tt.pageNumber)
			}
			if len(got.Lines) != tt.wantLines {
				t.Errorf("len(Page.Lines) = %d, want %d", len(got.Lines), tt.wantLines)
			}
		})
	}
}

func TestProcessPageHeadingType(t *testing.T) {
	// アウトラインが空の場合: ALL CAPS 行がフォールバックで h3 になる
	lines := []string{"CHAPTER ONE"}
	page := ProcessPage(lines, 1, nil)
	if len(page.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(page.Lines))
	}
	if page.Lines[0].Type != model.LineTypeHeading {
		t.Errorf("expected LineTypeHeading, got %v", page.Lines[0].Type)
	}
	if page.Lines[0].HeadingLevel != 3 {
		t.Errorf("expected heading level 3, got %d", page.Lines[0].HeadingLevel)
	}
}

func TestProcessPageOutlineHeading(t *testing.T) {
	// アウトラインマッチングで h1 が検出される
	lines := []string{"Chapter 1: Introduction", "Some body text."}
	outline := []model.OutlineItem{
		{Title: "Chapter 1: Introduction"},
	}
	page := ProcessPage(lines, 1, outline)
	if len(page.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(page.Lines))
	}
	if page.Lines[0].Type != model.LineTypeHeading {
		t.Errorf("expected LineTypeHeading for first line, got %v", page.Lines[0].Type)
	}
	if page.Lines[0].HeadingLevel != 1 {
		t.Errorf("expected heading level 1, got %d", page.Lines[0].HeadingLevel)
	}
	if page.Lines[1].Type != model.LineTypeNormal {
		t.Errorf("expected LineTypeNormal for second line, got %v", page.Lines[1].Type)
	}
}
