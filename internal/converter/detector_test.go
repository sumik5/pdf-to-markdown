package converter

import (
	"testing"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

func TestDetectHeadingsFromOutline(t *testing.T) {
	tests := []struct {
		name      string
		lines     []model.Line
		outline   []model.OutlineItem
		wantTypes []model.LineType
		wantLevels []int
	}{
		{
			name:      "no outline, normal text unchanged",
			lines:     []model.Line{{Text: "Hello world"}},
			outline:   []model.OutlineItem{{Title: "Introduction"}},
			wantTypes: []model.LineType{model.LineTypeNormal},
			wantLevels: []int{0},
		},
		{
			name:  "exact match - h1",
			lines: []model.Line{{Text: "Introduction"}},
			outline: []model.OutlineItem{
				{Title: "Introduction"},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{1},
		},
		{
			name:  "case insensitive match",
			lines: []model.Line{{Text: "introduction"}},
			outline: []model.OutlineItem{
				{Title: "Introduction"},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{1},
		},
		{
			name:  "contains match (page number appended)",
			lines: []model.Line{{Text: "Chapter 1: Background 5"}},
			outline: []model.OutlineItem{
				{Title: "Chapter 1: Background"},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{1},
		},
		{
			name:  "nested outline - depth 1 becomes h2",
			lines: []model.Line{{Text: "Section 1.1"}},
			outline: []model.OutlineItem{
				{
					Title: "Chapter 1",
					Children: []model.OutlineItem{
						{Title: "Section 1.1"},
					},
				},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{2},
		},
		{
			name:  "deep nesting capped at h3",
			lines: []model.Line{{Text: "Deep Section"}},
			outline: []model.OutlineItem{
				{
					Title: "Ch1",
					Children: []model.OutlineItem{
						{
							Title: "Sub1",
							Children: []model.OutlineItem{
								{
									Title: "Sub1.1",
									Children: []model.OutlineItem{
										{Title: "Deep Section"},
									},
								},
							},
						},
					},
				},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{3},
		},
		{
			name:  "no duplicate match (used flag)",
			lines: []model.Line{{Text: "Introduction"}, {Text: "Introduction"}},
			outline: []model.OutlineItem{
				{Title: "Introduction"},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading, model.LineTypeNormal},
			wantLevels: []int{1, 0},
		},
		{
			name:    "fallback ALL CAPS when outline empty",
			lines:   []model.Line{{Text: "CHAPTER ONE"}},
			outline: nil,
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{3},
		},
		{
			name:    "fallback ALL CAPS too long - not heading",
			lines:   []model.Line{{Text: "THIS IS A VERY LONG ALL CAPS LINE THAT EXCEEDS EIGHTY CHARACTERS IN TOTAL LENGTH HERE"}},
			outline: nil,
			wantTypes:  []model.LineType{model.LineTypeNormal},
			wantLevels: []int{0},
		},
		{
			name:    "no fallback when outline is non-empty but unmatched",
			lines:   []model.Line{{Text: "SOME CAPS LINE"}},
			outline: []model.OutlineItem{{Title: "Other Title"}},
			wantTypes:  []model.LineType{model.LineTypeNormal},
			wantLevels: []int{0},
		},
		{
			name:  "whitespace normalization in matching",
			lines: []model.Line{{Text: "Chapter  1:  Introduction"}},
			outline: []model.OutlineItem{
				{Title: "Chapter 1: Introduction"},
			},
			wantTypes:  []model.LineType{model.LineTypeHeading},
			wantLevels: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectHeadingsFromOutline(tt.lines, tt.outline)

			if len(got) != len(tt.wantTypes) {
				t.Fatalf("detectHeadingsFromOutline() = %d lines, want %d", len(got), len(tt.wantTypes))
			}

			for i, wantType := range tt.wantTypes {
				if got[i].Type != wantType {
					t.Errorf("line[%d].Type = %v, want %v (text: %q)", i, got[i].Type, wantType, tt.lines[i].Text)
				}
				if got[i].HeadingLevel != tt.wantLevels[i] {
					t.Errorf("line[%d].HeadingLevel = %d, want %d (text: %q)", i, got[i].HeadingLevel, tt.wantLevels[i], tt.lines[i].Text)
				}
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
			name:     "bullet with •",
			lines:    []model.Line{{Text: "• First item"}},
			wantType: []model.LineType{model.LineTypeBullet},
		},
		{
			name:     "bullet with dash",
			lines:    []model.Line{{Text: "- Second item"}},
			wantType: []model.LineType{model.LineTypeBullet},
		},
		{
			name:     "numbered list 1.",
			lines:    []model.Line{{Text: "1. First"}},
			wantType: []model.LineType{model.LineTypeNumbered},
		},
		{
			name:     "numbered list a.",
			lines:    []model.Line{{Text: "a. First"}},
			wantType: []model.LineType{model.LineTypeNumbered},
		},
		{
			name:     "normal text",
			lines:    []model.Line{{Text: "Normal paragraph text."}},
			wantType: []model.LineType{model.LineTypeNormal},
		},
		{
			name:     "heading not changed to list",
			lines:    []model.Line{{Text: "• Heading", Type: model.LineTypeHeading}},
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
