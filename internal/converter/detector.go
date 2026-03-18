package converter

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

var (
	bulletPattern   = regexp.MustCompile(`^[•·▪▫◦‣⁃]\s*|^[-–—]\s+|^\*\s+`)
	numberedPattern = regexp.MustCompile(`^\d+[.)]\s*|^[a-zA-Z][.)]\s*|^[ivxIVX]+[.)]\s*`)
)

// detectHeadings analyzes each line's font size and bold status to assign heading levels.
// Returns a new slice with Type and HeadingLevel set.
func detectHeadings(lines []model.Line) []model.Line {
	result := make([]model.Line, len(lines))
	for i, line := range lines {
		result[i] = line
		result[i].HeadingLevel = headingLevel(line)
		if result[i].HeadingLevel > 0 {
			result[i].Type = model.LineTypeHeading
		}
	}
	return result
}

// headingLevel returns 1/2/3 for headings, or 0 for normal text.
func headingLevel(line model.Line) int {
	if len(line.Items) == 0 {
		return 0
	}

	avgFontSize := averageFontSize(line.Items)
	isBold := anyBold(line.Items)
	isAllCaps := isAllUpperCase(line.Text)

	switch {
	case avgFontSize > 20 || (avgFontSize > 16 && isBold):
		return 1
	case avgFontSize > 16 || (avgFontSize > 14 && isBold):
		return 2
	case avgFontSize > 14 || (avgFontSize > 12 && isBold):
		return 3
	case isAllCaps && len(line.Text) < 50:
		return 3
	default:
		return 0
	}
}

// averageFontSize computes the mean font size across all items in a line.
func averageFontSize(items []model.TextItem) float64 {
	if len(items) == 0 {
		return 0
	}
	sum := 0.0
	for _, item := range items {
		sum += item.FontSize
	}
	return sum / float64(len(items))
}

// anyBold returns true if any item in the line has a bold font.
func anyBold(items []model.TextItem) bool {
	for _, item := range items {
		if item.IsBold() {
			return true
		}
	}
	return false
}

// isAllUpperCase returns true if the text contains only uppercase letters (and no lowercase).
func isAllUpperCase(text string) bool {
	hasLetter := false
	for _, r := range text {
		if unicode.IsLetter(r) {
			hasLetter = true
			if unicode.IsLower(r) {
				return false
			}
		}
	}
	return hasLetter
}

// detectLists examines line text to identify bullet and numbered list items.
func detectLists(lines []model.Line) []model.Line {
	result := make([]model.Line, len(lines))
	for i, line := range lines {
		result[i] = line
		if line.Type == model.LineTypeHeading {
			continue
		}
		if bulletPattern.MatchString(line.Text) {
			result[i].Type = model.LineTypeBullet
		} else if numberedPattern.MatchString(line.Text) {
			result[i].Type = model.LineTypeNumbered
		}
	}
	return result
}

// stripBulletPrefix removes bullet characters from the beginning of text.
func stripBulletPrefix(text string) string {
	stripped := bulletPattern.ReplaceAllString(text, "")
	return strings.TrimSpace(stripped)
}
