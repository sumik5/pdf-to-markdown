package extractor

import (
	"github.com/ledongthuc/pdf"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

// extractMetadata reads the PDF Info dictionary and returns a Metadata struct.
func extractMetadata(r *pdf.Reader) model.Metadata {
	info := r.Trailer().Key("Root").Key("Info")
	if info.IsNull() {
		// Some PDFs store Info directly under Trailer (not under Root)
		info = r.Trailer().Key("Info")
	}

	if info.IsNull() {
		return model.Metadata{}
	}

	return model.Metadata{
		Title:        info.Key("Title").Text(),
		Author:       info.Key("Author").Text(),
		Subject:      info.Key("Subject").Text(),
		Keywords:     info.Key("Keywords").Text(),
		Creator:      info.Key("Creator").Text(),
		Producer:     info.Key("Producer").Text(),
		CreationDate: info.Key("CreationDate").Text(),
		ModDate:      info.Key("ModDate").Text(),
	}
}
