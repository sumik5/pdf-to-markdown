package converter

import (
	"strings"

	"github.com/shivase/pdf-to-markdown/internal/model"
)

// parseLinesWithIndentation は生の行文字列スライスを model.Line に変換し、
// pdftotext -layout 出力の先頭スペース数からインデントレベルを算出する。
// 空行はスキップし、最小スペース数を基準(0)として4スペースごとに1レベルとする。
func parseLinesWithIndentation(rawLines []string) []model.Line {
	type entry struct {
		text          string
		leadingSpaces int
	}

	entries := make([]entry, 0, len(rawLines))
	for _, raw := range rawLines {
		trimmedRight := strings.TrimRight(raw, " \t\r\n")
		if trimmedRight == "" {
			continue
		}
		text := strings.Join(strings.Fields(strings.TrimSpace(trimmedRight)), " ")
		leadingSpaces := len(trimmedRight) - len(strings.TrimLeft(trimmedRight, " "))
		entries = append(entries, entry{text: text, leadingSpaces: leadingSpaces})
	}

	if len(entries) == 0 {
		return nil
	}

	// 全行の最小先頭スペース数を求め、基準(IndentLevel 0)とする
	minSpaces := entries[0].leadingSpaces
	for _, e := range entries[1:] {
		if e.leadingSpaces < minSpaces {
			minSpaces = e.leadingSpaces
		}
	}

	result := make([]model.Line, len(entries))
	for i, e := range entries {
		indentLevel := (e.leadingSpaces - minSpaces) / 4
		result[i] = model.Line{
			Text:        e.text,
			IndentLevel: indentLevel,
		}
	}
	return result
}
