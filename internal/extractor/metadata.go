package extractor

import (
	"os/exec"
	"strings"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

// extractMetadataWithPdfinfo runs pdfinfo and parses the output into a Metadata struct.
// If pdfinfo fails, an empty Metadata is returned (metadata is non-essential).
func extractMetadataWithPdfinfo(filePath string) (model.Metadata, error) {
	out, err := exec.Command("pdfinfo", filePath).Output()
	if err != nil {
		return model.Metadata{}, nil
	}

	var meta model.Metadata
	for _, line := range strings.Split(string(out), "\n") {
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		switch key {
		case "Title":
			meta.Title = value
		case "Author":
			meta.Author = value
		case "Subject":
			meta.Subject = value
		case "Keywords":
			meta.Keywords = value
		case "Creator":
			meta.Creator = value
		case "Producer":
			meta.Producer = value
		case "CreationDate":
			meta.CreationDate = value
		case "ModDate":
			meta.ModDate = value
		}
	}

	return meta, nil
}
