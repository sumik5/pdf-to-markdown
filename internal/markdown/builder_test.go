package markdown

import (
	"strings"
	"testing"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

func TestBuildEmpty(t *testing.T) {
	b := New()
	got := b.Build(nil, nil, nil)
	if got != "" {
		t.Errorf("Build() with no content = %q, want empty string", got)
	}
}

func TestBuildMetadata(t *testing.T) {
	b := New()
	meta := &model.Metadata{
		Title:  "Test Document",
		Author: "Author Name",
	}
	got := b.Build(nil, meta, nil)
	if !strings.Contains(got, "## Document Metadata") {
		t.Errorf("expected metadata header, got: %s", got)
	}
	if !strings.Contains(got, "**Title:** Test Document") {
		t.Errorf("expected title, got: %s", got)
	}
	if !strings.Contains(got, "**Author:** Author Name") {
		t.Errorf("expected author, got: %s", got)
	}
}

func TestBuildOutline(t *testing.T) {
	b := New()
	outline := []model.OutlineItem{
		{Title: "Chapter 1", Children: []model.OutlineItem{
			{Title: "Section 1.1"},
		}},
		{Title: "Chapter 2"},
	}
	got := b.Build(nil, nil, outline)
	if !strings.Contains(got, "## Table of Contents") {
		t.Errorf("expected TOC header, got: %s", got)
	}
	if !strings.Contains(got, "- Chapter 1") {
		t.Errorf("expected Chapter 1, got: %s", got)
	}
	if !strings.Contains(got, "  - Section 1.1") {
		t.Errorf("expected Section 1.1 indented, got: %s", got)
	}
}

func TestRenderLine(t *testing.T) {
	tests := []struct {
		name     string
		line     model.Line
		contains string
	}{
		{
			name: "h1 heading",
			line: model.Line{
				Type:         model.LineTypeHeading,
				HeadingLevel: 1,
				Text:         "Big Title",
			},
			contains: "# Big Title",
		},
		{
			name: "h2 heading",
			line: model.Line{
				Type:         model.LineTypeHeading,
				HeadingLevel: 2,
				Text:         "Subtitle",
			},
			contains: "## Subtitle",
		},
		{
			name: "h3 heading",
			line: model.Line{
				Type:         model.LineTypeHeading,
				HeadingLevel: 3,
				Text:         "Sub-subtitle",
			},
			contains: "### Sub-subtitle",
		},
		{
			name: "bullet list item",
			line: model.Line{
				Type: model.LineTypeBullet,
				Text: "• First item",
			},
			contains: "- First item",
		},
		{
			name: "numbered list item",
			line: model.Line{
				Type: model.LineTypeNumbered,
				Text: "1. First item",
			},
			contains: "1. First item",
		},
		{
			name: "normal text",
			line: model.Line{
				Type: model.LineTypeNormal,
				Text: "Paragraph text",
			},
			contains: "Paragraph text",
		},
		{
			name: "indented normal text",
			line: model.Line{
				Type:        model.LineTypeNormal,
				Text:        "Indented",
				IndentLevel: 2,
			},
			contains: "    Indented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderLine(tt.line)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("renderLine() = %q, want to contain %q", got, tt.contains)
			}
		})
	}
}

func TestBuildPagesWithSeparator(t *testing.T) {
	b := New()
	pages := []model.Page{
		{
			Number: 1,
			Lines: []model.Line{
				{Type: model.LineTypeNormal, Text: "Page 1 content"},
			},
		},
		{
			Number: 2,
			Lines: []model.Line{
				{Type: model.LineTypeNormal, Text: "Page 2 content"},
			},
		},
	}
	got := b.Build(pages, nil, nil)
	if !strings.Contains(got, "---") {
		t.Errorf("expected page separator ---, got: %s", got)
	}
	if !strings.Contains(got, "Page 1 content") {
		t.Errorf("expected page 1 content, got: %s", got)
	}
	if !strings.Contains(got, "Page 2 content") {
		t.Errorf("expected page 2 content, got: %s", got)
	}
}

func TestStripBulletPrefix(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"• First item", "First item"},
		{"· Second item", "Second item"},
		{"▪ Third item", "Third item"},
		{"- Dash item", "Dash item"},
		{"* Star item", "Star item"},
		{"Normal text", "Normal text"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stripBulletPrefix(tt.input)
			if got != tt.want {
				t.Errorf("stripBulletPrefix(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
