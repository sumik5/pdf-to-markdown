// Package model defines shared data types for the PDF to Markdown converter.
package model

import "strings"

// LineType represents the type of a text line.
type LineType int

const (
	// LineTypeNormal is a regular paragraph line.
	LineTypeNormal LineType = iota
	// LineTypeHeading is a heading line.
	LineTypeHeading
	// LineTypeBullet is a bullet list item.
	LineTypeBullet
	// LineTypeNumbered is a numbered list item.
	LineTypeNumbered
)

// TextItem represents an individual text element extracted from a PDF.
type TextItem struct {
	Text     string
	X        float64
	Y        float64
	Width    float64
	FontName string
	FontSize float64
}

// IsBold returns true if the font name indicates bold styling.
func (t TextItem) IsBold() bool {
	return strings.Contains(strings.ToLower(t.FontName), "bold")
}

// Line represents text items grouped into a single visual line.
type Line struct {
	Items        []TextItem
	Y            float64
	IndentLevel  int
	Type         LineType
	HeadingLevel int
	Text         string
}

// Page represents a single PDF page with its extracted lines.
type Page struct {
	Number int
	Lines  []Line
}

// Metadata holds PDF document metadata from the Info dictionary.
type Metadata struct {
	Title        string
	Author       string
	Subject      string
	Keywords     string
	Creator      string
	Producer     string
	CreationDate string
	ModDate      string
}

// IsEmpty returns true if no metadata fields are populated.
func (m *Metadata) IsEmpty() bool {
	return m.Title == "" && m.Author == "" && m.Subject == "" &&
		m.Keywords == "" && m.Creator == "" && m.Producer == "" &&
		m.CreationDate == "" && m.ModDate == ""
}

// OutlineItem represents a bookmark entry in the PDF outline tree.
type OutlineItem struct {
	Title    string
	Children []OutlineItem
}

// ConversionResult holds statistics about the conversion operation.
type ConversionResult struct {
	Pages      int    `json:"pages"`
	Characters int    `json:"characters"`
	OutputFile string `json:"outputFile"`
}
