// Package extractor handles reading PDF files and extracting structured content.
package extractor

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

// PDFExtractor extracts text, metadata, and outline from a PDF file.
type PDFExtractor struct {
	filePath string
}

// New creates a new PDFExtractor for the given file path.
func New(filePath string) *PDFExtractor {
	return &PDFExtractor{filePath: filePath}
}

// Extract opens the PDF and returns pages, metadata, outline, and page count.
func (e *PDFExtractor) Extract() (pages [][]model.TextItem, meta *model.Metadata, outline []model.OutlineItem, numPages int, err error) {
	f, r, err := pdf.Open(e.filePath)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("open PDF: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close PDF: %w", cerr)
		}
	}()

	numPages = r.NumPage()

	pages = make([][]model.TextItem, 0, numPages)
	for i := 1; i <= numPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			pages = append(pages, nil)
			continue
		}

		items := extractPageItems(p)
		pages = append(pages, items)
	}

	metaVal := extractMetadata(r)
	meta = &metaVal
	outline = extractOutline(r)

	return pages, meta, outline, numPages, nil
}

// extractPageItems converts a PDF page's raw text elements to TextItems.
func extractPageItems(p pdf.Page) []model.TextItem {
	content := p.Content()
	items := make([]model.TextItem, 0, len(content.Text))

	for _, t := range content.Text {
		text := strings.TrimSpace(t.S)
		if text == "" {
			continue
		}
		items = append(items, model.TextItem{
			Text:     text,
			X:        t.X,
			Y:        t.Y,
			Width:    t.W,
			FontName: t.Font,
			FontSize: t.FontSize,
		})
	}

	return items
}
