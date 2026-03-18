// pdf-to-markdown converts PDF files to Markdown format.
// Usage: pdf-to-markdown <input.pdf> <output.md>
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shivase/pdf-to-markdown/internal/converter"
	"github.com/shivase/pdf-to-markdown/internal/extractor"
	"github.com/shivase/pdf-to-markdown/internal/markdown"
	"github.com/shivase/pdf-to-markdown/internal/model"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: pdf-to-markdown <input.pdf> <output.md>")
	}

	inputPath := args[0]
	outputPath := args[1]

	// Extract content from PDF
	ext := extractor.New(inputPath)
	pages, meta, outline, numPages, err := ext.Extract()
	if err != nil {
		return fmt.Errorf("extract PDF: %w", err)
	}

	// Convert raw text items to structured lines
	structuredPages := make([]model.Page, 0, len(pages))
	for i, lines := range pages {
		page := converter.ProcessPage(lines, i+1, outline)
		structuredPages = append(structuredPages, page)
	}

	// Build Markdown output
	builder := markdown.New()
	content := builder.Build(structuredPages, meta, outline)

	// Write output file
	if err := os.WriteFile(outputPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	// Output conversion stats to stderr as JSON
	result := model.ConversionResult{
		Pages:      numPages,
		Characters: len(content),
		OutputFile: outputPath,
	}
	if err := json.NewEncoder(os.Stderr).Encode(result); err != nil {
		return fmt.Errorf("encode result: %w", err)
	}

	fmt.Printf("変換完了: %s\n", outputPath)

	return nil
}
