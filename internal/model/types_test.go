package model

import "testing"

func TestIsBold(t *testing.T) {
	tests := []struct {
		name     string
		fontName string
		want     bool
	}{
		{"bold font", "Helvetica-Bold", true},
		{"bold lower", "times-bold", true},
		{"mixed bold", "Arial-BoldItalic", true},
		{"regular font", "Helvetica", false},
		{"empty font", "", false},
		{"courier", "Courier", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := TextItem{FontName: tt.fontName}
			got := item.IsBold()
			if got != tt.want {
				t.Errorf("IsBold() = %v, want %v (font: %q)", got, tt.want, tt.fontName)
			}
		})
	}
}

func TestMetadataIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		meta *Metadata
		want bool
	}{
		{"empty metadata", &Metadata{}, true},
		{"title only", &Metadata{Title: "Test"}, false},
		{"author only", &Metadata{Author: "Author"}, false},
		{"producer only", &Metadata{Producer: "PDF Generator"}, false},
		{"full metadata", &Metadata{Title: "T", Author: "A", Subject: "S"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.meta.IsEmpty()
			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
