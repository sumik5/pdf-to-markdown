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

func TestCheckCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "存在するコマンド pdftotext",
			command:     "pdftotext",
			expectError: false,
		},
		{
			name:        "存在しないコマンド",
			command:     "this-command-does-not-exist-xyzzy",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCommand(tt.command)
			if tt.expectError && err == nil {
				t.Errorf("checkCommand(%q) expected error, got nil", tt.command)
			}
			if !tt.expectError && err != nil {
				t.Errorf("checkCommand(%q) unexpected error: %v", tt.command, err)
			}
		})
	}
}

func TestExtractMetadataWithPdfinfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantMeta func(t *testing.T, got interface{})
	}{
		{
			name:  "存在しないファイルは空のMetadataを返す",
			input: "/non/existent/file.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := extractMetadataWithPdfinfo(tt.input)
			// pdfinfoが失敗してもエラーは返さない仕様
			if err != nil {
				t.Errorf("extractMetadataWithPdfinfo() should not return error, got: %v", err)
			}
			// 存在しないファイルの場合は空のMetadataが返る
			if meta.Title != "" || meta.Author != "" {
				t.Errorf("expected empty Metadata for non-existent file, got: %+v", meta)
			}
		})
	}
}

func TestExtractWithPdftotext_InvalidFile(t *testing.T) {
	_, _, err := extractWithPdftotext("/non/existent/file.pdf")
	if err == nil {
		t.Error("extractWithPdftotext() should return error for non-existent file")
	}
}
