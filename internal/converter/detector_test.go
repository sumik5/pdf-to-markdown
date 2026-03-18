package converter

import (
	"testing"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

func TestHeadingLevel(t *testing.T) {
	tests := []struct {
		name string
		line model.Line
		want int
	}{
		{
			name: "empty line",
			line: model.Line{},
			want: 0,
		},
		{
			name: "normal text (12pt)",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 12, FontName: "Helvetica"}},
				Text:  "Normal text here",
			},
			want: 0,
		},
		{
			name: "h1 - large font (>20pt)",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 24, FontName: "Times-Roman"}},
				Text:  "Big heading",
			},
			want: 1,
		},
		{
			name: "h1 - bold 17pt",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 17, FontName: "Helvetica-Bold"}},
				Text:  "Bold heading",
			},
			want: 1,
		},
		{
			name: "h2 - 17pt regular",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 17, FontName: "Times-Roman"}},
				Text:  "Subheading",
			},
			want: 2,
		},
		{
			name: "h2 - bold 15pt",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 15, FontName: "Times-Bold"}},
				Text:  "Subheading bold",
			},
			want: 2,
		},
		{
			name: "h3 - 15pt regular",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 15, FontName: "Times-Roman"}},
				Text:  "Small heading",
			},
			want: 3,
		},
		{
			name: "h3 - bold 13pt",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 13, FontName: "Helvetica-Bold"}},
				Text:  "Small bold",
			},
			want: 3,
		},
		{
			name: "h3 - all caps short",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 12, FontName: "Helvetica"}},
				Text:  "INTRODUCTION",
			},
			want: 3,
		},
		{
			name: "not h3 - all caps long",
			line: model.Line{
				Items: []model.TextItem{{FontSize: 12, FontName: "Helvetica"}},
				Text:  "THIS IS A VERY LONG ALL CAPS TEXT THAT EXCEEDS FIFTY CHARACTERS LIMIT",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := headingLevel(tt.line)
			if got != tt.want {
				t.Errorf("headingLevel() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestIsAllUpperCase(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"all uppercase", "HELLO WORLD", true},
		{"mixed case", "Hello World", false},
		{"all lowercase", "hello world", false},
		{"uppercase with numbers", "CHAPTER 1", true},
		{"empty string", "", false},
		{"only numbers", "123", false},
		{"uppercase with symbols", "TITLE: PART ONE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAllUpperCase(tt.text)
			if got != tt.want {
				t.Errorf("isAllUpperCase(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestDetectLists(t *testing.T) {
	tests := []struct {
		name     string
		lines    []model.Line
		wantType []model.LineType
	}{
		{
			name: "bullet with •",
			lines: []model.Line{
				{Text: "• First item", Items: []model.TextItem{{FontSize: 12}}},
			},
			wantType: []model.LineType{model.LineTypeBullet},
		},
		{
			name: "bullet with dash",
			lines: []model.Line{
				{Text: "- Second item", Items: []model.TextItem{{FontSize: 12}}},
			},
			wantType: []model.LineType{model.LineTypeBullet},
		},
		{
			name: "numbered list 1.",
			lines: []model.Line{
				{Text: "1. First", Items: []model.TextItem{{FontSize: 12}}},
			},
			wantType: []model.LineType{model.LineTypeNumbered},
		},
		{
			name: "numbered list a.",
			lines: []model.Line{
				{Text: "a. First", Items: []model.TextItem{{FontSize: 12}}},
			},
			wantType: []model.LineType{model.LineTypeNumbered},
		},
		{
			name: "normal text",
			lines: []model.Line{
				{Text: "Normal paragraph text.", Items: []model.TextItem{{FontSize: 12}}},
			},
			wantType: []model.LineType{model.LineTypeNormal},
		},
		{
			name: "heading not changed to list",
			lines: []model.Line{
				{Text: "• Heading", Type: model.LineTypeHeading, Items: []model.TextItem{{FontSize: 12}}},
			},
			wantType: []model.LineType{model.LineTypeHeading},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectLists(tt.lines)
			for i, wantType := range tt.wantType {
				if got[i].Type != wantType {
					t.Errorf("line[%d].Type = %v, want %v (text: %q)", i, got[i].Type, wantType, tt.lines[i].Text)
				}
			}
		})
	}
}

func TestStripBulletPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bullet •", "• First item", "First item"},
		{"bullet -", "- Second item", "Second item"},
		{"bullet *", "* Third item", "Third item"},
		{"no bullet", "Normal text", "Normal text"},
		{"bullet ▪", "▪ item", "item"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripBulletPrefix(tt.input)
			if got != tt.want {
				t.Errorf("stripBulletPrefix(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAverageFontSize(t *testing.T) {
	tests := []struct {
		name  string
		items []model.TextItem
		want  float64
	}{
		{"empty", nil, 0},
		{"single", []model.TextItem{{FontSize: 12}}, 12},
		{"average", []model.TextItem{{FontSize: 10}, {FontSize: 20}}, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := averageFontSize(tt.items)
			if got != tt.want {
				t.Errorf("averageFontSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
