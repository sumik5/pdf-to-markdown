// Package markdown builds Markdown strings from structured PDF content.
package markdown

import (
	"fmt"
	"strings"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

// Builder constructs a Markdown document from PDF pages, metadata, and outline.
type Builder struct {
	sb strings.Builder
}

// New creates a new Builder.
func New() *Builder {
	return &Builder{}
}

// Build constructs and returns the full Markdown document.
func (b *Builder) Build(pages []model.Page, meta *model.Metadata, outline []model.OutlineItem) string {
	b.sb.Reset()

	if meta != nil && !meta.IsEmpty() {
		b.writeMetadata(meta)
		b.sb.WriteString("\n\n")
	}

	if len(outline) > 0 {
		b.writeOutline(outline)
		b.sb.WriteString("\n\n")
	}

	pageContents := make([]string, 0, len(pages))
	for _, page := range pages {
		content := b.buildPage(page)
		if strings.TrimSpace(content) != "" {
			pageContents = append(pageContents, content)
		}
	}

	b.sb.WriteString(strings.Join(pageContents, "\n\n---\n\n"))

	return b.sb.String()
}

// buildPage converts a single Page's lines to Markdown text.
func (b *Builder) buildPage(page model.Page) string {
	var sb strings.Builder

	for _, line := range page.Lines {
		if line.Text == "" {
			continue
		}
		sb.WriteString(renderLine(line))
	}

	return strings.TrimSpace(sb.String())
}

// renderLine converts a single Line to its Markdown representation.
func renderLine(line model.Line) string {
	indent := strings.Repeat("  ", line.IndentLevel)

	switch line.Type {
	case model.LineTypeHeading:
		prefix := strings.Repeat("#", line.HeadingLevel)
		return fmt.Sprintf("%s %s\n\n", prefix, line.Text)
	case model.LineTypeBullet:
		text := stripBulletPrefix(line.Text)
		return fmt.Sprintf("%s- %s\n", indent, text)
	case model.LineTypeNumbered:
		return fmt.Sprintf("%s%s\n", indent, line.Text)
	case model.LineTypeNormal:
		if line.IndentLevel > 0 {
			return fmt.Sprintf("%s%s\n", indent, line.Text)
		}
		return fmt.Sprintf("%s\n", line.Text)
	}
	return ""
}

// stripBulletPrefix removes common bullet symbols from text.
func stripBulletPrefix(text string) string {
	// Remove leading bullet characters and whitespace
	if text == "" {
		return text
	}

	bulletChars := "•·▪▫◦‣⁃"
	dashPrefixes := []string{"- ", "– ", "— ", "* "}

	r := []rune(text)
	if strings.ContainsRune(bulletChars, r[0]) {
		text = strings.TrimSpace(string(r[1:]))
		return text
	}

	for _, prefix := range dashPrefixes {
		if strings.HasPrefix(text, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(text, prefix))
		}
	}

	return text
}

// writeMetadata renders the PDF metadata section.
func (b *Builder) writeMetadata(meta *model.Metadata) {
	b.sb.WriteString("## Document Metadata\n\n")
	if meta.Title != "" {
		fmt.Fprintf(&b.sb, "**Title:** %s\n", meta.Title)
	}
	if meta.Author != "" {
		fmt.Fprintf(&b.sb, "**Author:** %s\n", meta.Author)
	}
	if meta.Subject != "" {
		fmt.Fprintf(&b.sb, "**Subject:** %s\n", meta.Subject)
	}
	if meta.Keywords != "" {
		fmt.Fprintf(&b.sb, "**Keywords:** %s\n", meta.Keywords)
	}
	if meta.Creator != "" {
		fmt.Fprintf(&b.sb, "**Creator:** %s\n", meta.Creator)
	}
	if meta.Producer != "" {
		fmt.Fprintf(&b.sb, "**Producer:** %s\n", meta.Producer)
	}
	if meta.CreationDate != "" {
		fmt.Fprintf(&b.sb, "**Created:** %s\n", meta.CreationDate)
	}
	if meta.ModDate != "" {
		fmt.Fprintf(&b.sb, "**Modified:** %s\n", meta.ModDate)
	}
}

// writeOutline renders the PDF outline (table of contents) section.
func (b *Builder) writeOutline(outline []model.OutlineItem) {
	b.sb.WriteString("## Table of Contents\n\n")
	b.writeOutlineItems(outline, 0)
}

// writeOutlineItems recursively renders outline items.
func (b *Builder) writeOutlineItems(items []model.OutlineItem, level int) {
	indent := strings.Repeat("  ", level)
	for _, item := range items {
		title := item.Title
		if title == "" {
			title = "Untitled"
		}
		fmt.Fprintf(&b.sb, "%s- %s\n", indent, title)
		if len(item.Children) > 0 {
			b.writeOutlineItems(item.Children, level+1)
		}
	}
}
