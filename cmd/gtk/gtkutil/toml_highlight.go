//go:build gtk

package gtkutil

import (
	"strings"
	"unicode"
)

// TOML highlight tag names (shared with GTK tag table).
const (
	TOMLTagComment = "toml-comment"
	TOMLTagString  = "toml-string"
	TOMLTagKey     = "toml-key"
	TOMLTagNumber  = "toml-number"
	TOMLTagBoolean = "toml-boolean"
	TOMLTagTable   = "toml-table"
)

// TOMLHighlightRange is a half-open byte interval [Start, End) in the full document.
type TOMLHighlightRange struct {
	Start int
	End   int
	Tag   string
}

// ComputeTOMLHighlightRanges returns syntax highlight spans for a TOML document.
func ComputeTOMLHighlightRanges(content string) []TOMLHighlightRange {
	lines := strings.Split(content, "\n")
	ranges := make([]TOMLHighlightRange, 0, 2)
	lineOffset := 0

	for _, line := range lines {
		ranges = append(ranges, computeTOMLLineHighlightRanges(line, lineOffset)...)
		lineOffset += len(line) + 1
	}

	return ranges
}

func computeTOMLLineHighlightRanges(line string, lineOffset int) []TOMLHighlightRange {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil
	}

	if strings.HasPrefix(trimmed, "#") {
		hashIdx := strings.Index(line, "#")
		if hashIdx < 0 {
			hashIdx = 0
		}

		return []TOMLHighlightRange{{
			Start: lineOffset + hashIdx,
			End:   lineOffset + len(line),
			Tag:   TOMLTagComment,
		}}
	}

	if strings.HasPrefix(trimmed, "[") {
		bracket := strings.Index(line, "[")
		if bracket < 0 {
			bracket = 0
		}

		return []TOMLHighlightRange{{
			Start: lineOffset + bracket,
			End:   lineOffset + len(line),
			Tag:   TOMLTagTable,
		}}
	}

	var ranges []TOMLHighlightRange

	commentIdx := findTOMLCommentIndex(line)
	lineContent := line
	if commentIdx >= 0 {
		ranges = append(ranges, TOMLHighlightRange{
			Start: lineOffset + commentIdx,
			End:   lineOffset + len(line),
			Tag:   TOMLTagComment,
		})
		lineContent = line[:commentIdx]
	}

	eqIdx := strings.Index(lineContent, "=")
	if eqIdx < 0 {
		return ranges
	}

	keyPart := strings.TrimSpace(lineContent[:eqIdx])
	if keyPart != "" {
		keyStart := strings.Index(lineContent, keyPart)
		if keyStart >= 0 {
			ranges = append(ranges, TOMLHighlightRange{
				Start: lineOffset + keyStart,
				End:   lineOffset + keyStart + len(keyPart),
				Tag:   TOMLTagKey,
			})
		}
	}

	valuePart := strings.TrimSpace(lineContent[eqIdx+1:])
	if valuePart == "" {
		return ranges
	}

	valueStart := eqIdx + 1
	for valueStart < len(lineContent) && (lineContent[valueStart] == ' ' || lineContent[valueStart] == '\t') {
		valueStart++
	}

	ranges = append(ranges, computeTOMLValueHighlightRanges(lineContent, lineOffset, valueStart, valuePart)...)

	return ranges
}

func computeTOMLValueHighlightRanges(lineContent string, lineOffset,
	valueStart int, value string,
) []TOMLHighlightRange {
	absValueStart := lineOffset + valueStart

	switch {
	case strings.HasPrefix(value, `"`) || strings.HasPrefix(value, `'`):
		return []TOMLHighlightRange{{
			Start: absValueStart,
			End:   lineOffset + len(lineContent),
			Tag:   TOMLTagString,
		}}
	case value == "true" || value == "false":
		return []TOMLHighlightRange{{
			Start: absValueStart,
			End:   lineOffset + len(lineContent),
			Tag:   TOMLTagBoolean,
		}}
	default:
		token := value
		if i := strings.IndexAny(token, " \t#"); i >= 0 {
			token = token[:i]
		}
		if isTOMLNumber(token) {
			return []TOMLHighlightRange{{
				Start: absValueStart,
				End:   absValueStart + len(token),
				Tag:   TOMLTagNumber,
			}}
		}
	}

	return nil
}

func findTOMLCommentIndex(line string) int {
	inSingle := false
	inDouble := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		switch ch {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return i
			}
		}
	}

	return -1
}

func isTOMLNumber(str string) bool {
	if str == "" {
		return false
	}

	i := 0
	if str[0] == '-' || str[0] == '+' {
		i++
	}
	if i >= len(str) {
		return false
	}

	hasDigit := false
	for ; i < len(str); i++ {
		chr := str[i]
		if unicode.IsDigit(rune(chr)) {
			hasDigit = true

			continue
		}
		if chr == '.' || chr == '_' || chr == 'e' || chr == 'E' || chr == '-' || chr == '+' {
			continue
		}

		return false
	}

	return hasDigit
}
