package converter

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

var (
	bulletPattern   = regexp.MustCompile(`^[•·▪▫◦‣⁃]\s*|^[-–—]\s+|^\*\s+`)
	numberedPattern = regexp.MustCompile(`^\d+[.)]\s*|^[a-zA-Z][.)]\s*|^[ivxIVX]+[.)]\s*`)
)

// outlineEntry はフラット化されたアウトラインエントリを表す。
type outlineEntry struct {
	title string
	level int // 1-based heading level
}

// flattenOutline はネストされた OutlineItem を再帰的にフラット化する。
// depth 0 → HeadingLevel 1, depth 1 → HeadingLevel 2, depth 2+ → HeadingLevel 3
func flattenOutline(items []model.OutlineItem, depth int) []outlineEntry {
	var result []outlineEntry
	level := depth + 1
	if level > 3 {
		level = 3
	}
	for _, item := range items {
		result = append(result, outlineEntry{title: item.Title, level: level})
		result = append(result, flattenOutline(item.Children, depth+1)...)
	}
	return result
}

// normalizeText は空白を正規化して小文字に変換する。
// pdftotext出力では単語間スペースが複数になる場合があるため、Fields で分割して再結合する。
func normalizeText(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// detectHeadingsFromOutline はアウトライン（目次）情報と行テキストを照合して見出しを検出する。
// マッチングは完全一致 → 含有一致（正規化後）の順で試みる。
// アウトラインが空の場合はフォールバックとして ALL CAPS ヒューリスティックを使用する。
func detectHeadingsFromOutline(lines []model.Line, outline []model.OutlineItem) []model.Line {
	entries := flattenOutline(outline, 0)
	used := make([]bool, len(entries))

	result := make([]model.Line, len(lines))
	for i, line := range lines {
		result[i] = line

		matched := false
		normalizedLine := normalizeText(line.Text)

		for j, entry := range entries {
			if used[j] {
				continue
			}
			normalizedTitle := normalizeText(entry.title)

			if strings.EqualFold(line.Text, entry.title) ||
				strings.Contains(normalizedLine, normalizedTitle) {
				result[i].Type = model.LineTypeHeading
				result[i].HeadingLevel = entry.level
				used[j] = true
				matched = true
				break
			}
		}

		// アウトラインが空の場合のフォールバック: ALL CAPS かつ短い行を h3 とみなす
		if !matched && len(outline) == 0 {
			if isAllUpperCase(line.Text) && len(line.Text) < 80 {
				result[i].Type = model.LineTypeHeading
				result[i].HeadingLevel = 3
			}
		}
	}
	return result
}

// isAllUpperCase returns true if the text contains only uppercase letters (and no lowercase).
func isAllUpperCase(text string) bool {
	hasLetter := false
	for _, r := range text {
		if unicode.IsLetter(r) {
			hasLetter = true
			if unicode.IsLower(r) {
				return false
			}
		}
	}
	return hasLetter
}

// detectLists examines line text to identify bullet and numbered list items.
func detectLists(lines []model.Line) []model.Line {
	result := make([]model.Line, len(lines))
	for i, line := range lines {
		result[i] = line
		if line.Type == model.LineTypeHeading {
			continue
		}
		if bulletPattern.MatchString(line.Text) {
			result[i].Type = model.LineTypeBullet
		} else if numberedPattern.MatchString(line.Text) {
			result[i].Type = model.LineTypeNumbered
		}
	}
	return result
}

// stripBulletPrefix removes bullet characters from the beginning of text.
func stripBulletPrefix(text string) string {
	stripped := bulletPattern.ReplaceAllString(text, "")
	return strings.TrimSpace(stripped)
}
