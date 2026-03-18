package converter

import (
	"math"
	"sort"
	"strings"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

const yTolerance = 2.0

// groupIntoLines groups TextItems that share approximately the same Y coordinate
// into visual lines, then sorts each line's items by X (left to right).
func groupIntoLines(items []model.TextItem) []model.Line {
	// Map from representative Y → items
	type lineEntry struct {
		y     float64
		items []model.TextItem
	}

	entries := make([]lineEntry, 0)
	yIndex := make(map[int]int) // rounded Y → index in entries

	for _, item := range items {
		roundedY := int(math.Round(item.Y))
		foundIdx := -1

		// Search nearby Y values within tolerance
		for ry, idx := range yIndex {
			if math.Abs(float64(ry)-item.Y) < yTolerance {
				foundIdx = idx
				break
			}
		}

		if foundIdx == -1 {
			yIndex[roundedY] = len(entries)
			entries = append(entries, lineEntry{y: item.Y, items: []model.TextItem{item}})
		} else {
			entries[foundIdx].items = append(entries[foundIdx].items, item)
		}
	}

	// Sort lines top to bottom (PDF Y increases bottom→top, so descending Y = top→bottom)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].y > entries[j].y
	})

	lines := make([]model.Line, 0, len(entries))
	for _, entry := range entries {
		// Sort items within line left to right
		sort.Slice(entry.items, func(i, j int) bool {
			return entry.items[i].X < entry.items[j].X
		})

		text := joinItemTexts(entry.items)
		lines = append(lines, model.Line{
			Items: entry.items,
			Y:     entry.y,
			Text:  text,
		})
	}

	return lines
}

// joinItemTexts concatenates TextItem texts, adding a space between items
// that are not adjacent (gap between end of previous and start of next).
func joinItemTexts(items []model.TextItem) string {
	if len(items) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(items[0].Text)

	for i := 1; i < len(items); i++ {
		prev := items[i-1]
		curr := items[i]
		gap := curr.X - (prev.X + prev.Width)
		// If gap is more than 1 point, insert a space
		if gap > 1.0 || !strings.HasSuffix(prev.Text, " ") {
			sb.WriteString(" ")
		}
		sb.WriteString(curr.Text)
	}

	return strings.TrimSpace(sb.String())
}

// detectIndentation calculates the indent level for each line based on its
// X offset from the minimum X across all lines.
func detectIndentation(lines []model.Line) []model.Line {
	minX := findMinX(lines)
	result := make([]model.Line, len(lines))

	for i, line := range lines {
		firstX := 0.0
		if len(line.Items) > 0 {
			firstX = line.Items[0].X
		}
		indent := int(math.Round((firstX - minX) / 10.0))
		if indent < 0 {
			indent = 0
		}
		result[i] = line
		result[i].IndentLevel = indent
	}

	return result
}

// findMinX returns the minimum X coordinate across all lines (only positive values).
func findMinX(lines []model.Line) float64 {
	minX := math.MaxFloat64
	for _, line := range lines {
		if len(line.Items) > 0 && line.Items[0].X > 0 {
			if line.Items[0].X < minX {
				minX = line.Items[0].X
			}
		}
	}
	if minX == math.MaxFloat64 {
		return 0
	}
	return minX
}
