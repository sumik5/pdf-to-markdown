// Package converter orchestrates the transformation of raw PDF text lines into
// structured Lines with heading, list, and indentation information.
package converter

import (
	"github.com/shivase/pdf-to-markdown/internal/model"
)

// ProcessPage converts a slice of raw text lines (from one PDF page) into
// structured Lines ready for Markdown rendering.
func ProcessPage(lines []string, pageNumber int, outline []model.OutlineItem) model.Page {
	parsedLines := parseLinesWithIndentation(lines)
	parsedLines = detectHeadingsFromOutline(parsedLines, outline)
	parsedLines = detectLists(parsedLines)

	return model.Page{
		Number: pageNumber,
		Lines:  parsedLines,
	}
}
