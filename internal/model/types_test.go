package model

import "testing"

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
