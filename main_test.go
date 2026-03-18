package main

import (
	"os"
	"testing"
)

func TestRunNoArgs(t *testing.T) {
	err := run([]string{})
	if err == nil {
		t.Error("run() with no args should return error")
	}
}

func TestRunOneArg(t *testing.T) {
	err := run([]string{"input.pdf"})
	if err == nil {
		t.Error("run() with one arg should return error")
	}
}

func TestRunNonExistentPDF(t *testing.T) {
	err := run([]string{"/non/existent/file.pdf", "/tmp/output.md"})
	if err == nil {
		t.Error("run() with non-existent PDF should return error")
	}
}

func TestRunInvalidPDF(t *testing.T) {
	// Create a temp file that is not a valid PDF
	tmp, err := os.CreateTemp("", "invalid*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString("not a pdf file content"); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	runErr := run([]string{tmp.Name(), "/tmp/output_invalid.md"})
	if runErr == nil {
		t.Error("run() with invalid PDF should return error")
		os.Remove("/tmp/output_invalid.md")
	}
}
