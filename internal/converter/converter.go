// Package converter orchestrates the transformation of raw PDF text items into
// structured Lines with heading, list, and indentation information.
package converter

import (
	"github.com/shivase/pdf-to-markdown/internal/model"
)

// ProcessPage converts a slice of raw TextItems (from one PDF page) into
// structured Lines ready for Markdown rendering.
func ProcessPage(items []model.TextItem, pageNumber int) model.Page {
	lines := groupIntoLines(items)
	lines = detectIndentation(lines)
	lines = detectHeadings(lines)
	lines = detectLists(lines)

	return model.Page{
		Number: pageNumber,
		Lines:  lines,
	}
}
