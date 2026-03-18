package extractor

import (
	"testing"
)

func TestNew(t *testing.T) {
	e := New("test.pdf")
	if e == nil {
		t.Fatal("New() should return non-nil PDFExtractor")
	}
	if e.filePath != "test.pdf" {
		t.Errorf("expected filePath %q, got %q", "test.pdf", e.filePath)
	}
}

func TestExtractNonExistentFile(t *testing.T) {
	e := New("/non/existent/file.pdf")
	_, _, _, _, err := e.Extract() //nolint:dogsled // 5-value return requires blank identifiers
	if err == nil {
		t.Error("Extract() should return error for non-existent file")
	}
}
