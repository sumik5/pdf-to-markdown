package converter

import (
	"testing"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

func TestProcessPage(t *testing.T) {
	tests := []struct {
		name       string
		items      []model.TextItem
		pageNumber int
		wantLines  int
	}{
		{
			name:       "empty page",
			items:      nil,
			pageNumber: 1,
			wantLines:  0,
		},
		{
			name: "single line page",
			items: []model.TextItem{
				{Text: "Hello World", X: 10, Y: 100, FontSize: 12},
			},
			pageNumber: 1,
			wantLines:  1,
		},
		{
			name: "heading detected",
			items: []model.TextItem{
				{Text: "Big Title", X: 10, Y: 100, FontSize: 24},
				{Text: "Normal text here.", X: 10, Y: 80, FontSize: 12},
			},
			pageNumber: 1,
			wantLines:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessPage(tt.items, tt.pageNumber)
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
	items := []model.TextItem{
		{Text: "CHAPTER ONE", X: 10, Y: 100, FontSize: 12, FontName: "Helvetica"},
	}
	page := ProcessPage(items, 1)
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
