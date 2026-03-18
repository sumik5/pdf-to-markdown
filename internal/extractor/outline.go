package extractor

import (
	"github.com/ledongthuc/pdf"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

// extractOutline reads the PDF Outlines (bookmark) tree and returns a nested slice.
func extractOutline(r *pdf.Reader) []model.OutlineItem {
	root := r.Trailer().Key("Root")
	if root.IsNull() {
		return nil
	}

	outlines := root.Key("Outlines")
	if outlines.IsNull() {
		return nil
	}

	first := outlines.Key("First")
	if first.IsNull() {
		return nil
	}

	return traverseOutline(first)
}

// traverseOutline recursively walks the linked-list structure of PDF outlines.
func traverseOutline(node pdf.Value) []model.OutlineItem {
	var items []model.OutlineItem

	for !node.IsNull() {
		title := node.Key("Title").Text()
		item := model.OutlineItem{Title: title}

		first := node.Key("First")
		if !first.IsNull() {
			item.Children = traverseOutline(first)
		}

		items = append(items, item)
		node = node.Key("Next")
	}

	return items
}
