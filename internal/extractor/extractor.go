// Package extractor handles reading PDF files and extracting structured content.
package extractor

import (
	"fmt"
	"os/exec"
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
func (e *PDFExtractor) Extract() (pages [][]string, meta *model.Metadata, outline []model.OutlineItem, numPages int, err error) {
	if err := checkCommand("pdftotext"); err != nil {
		return nil, nil, nil, 0, err
	}

	pages, numPages, err = extractWithPdftotext(e.filePath)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	metaVal, _ := extractMetadataWithPdfinfo(e.filePath)
	meta = &metaVal

	f, r, err := pdf.Open(e.filePath)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("open PDF for outline: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close PDF: %w", cerr)
		}
	}()

	outline = extractOutline(r)

	return pages, meta, outline, numPages, nil
}

// checkCommand verifies that the given command is available in PATH.
func checkCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("コマンド %q が見つかりません。poppler-utils をインストールしてください (brew install poppler / apt install poppler-utils): %w", name, err)
	}
	return nil
}

// extractWithPdftotext runs pdftotext -layout and splits output into pages and lines.
func extractWithPdftotext(filePath string) ([][]string, int, error) {
	out, err := exec.Command("pdftotext", "-layout", filePath, "-").Output()
	if err != nil {
		return nil, 0, fmt.Errorf("pdftotext の実行に失敗しました: %w", err)
	}

	raw := string(out)
	// form feed (\f) でページ分割
	pageTexts := strings.Split(raw, "\f")

	// pdftotext は末尾に余分な \f を付けることがある — 末尾の空ページを除去
	if len(pageTexts) > 0 && strings.TrimSpace(pageTexts[len(pageTexts)-1]) == "" {
		pageTexts = pageTexts[:len(pageTexts)-1]
	}

	pages := make([][]string, len(pageTexts))
	for i, pageText := range pageTexts {
		pages[i] = strings.Split(pageText, "\n")
	}

	return pages, len(pages), nil
}
