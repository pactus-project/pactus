// Package moniker reads the moniker of a node.
//
// A moniker is a free string that a node declares about itself. It can be a
// plain name like "John Doe", or it can follow the structured format of PIP-45:
//
//	name=john_doe;email=john.doe@example.com;web=https://example.com
//
// A moniker is not verified by anyone: any node can declare any name, email
// or website. The values must be shown as text that the node itself declared,
// and never as a verified identity.
package moniker

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// MaxLength is the maximum length of a moniker, in bytes.
	MaxLength = 128

	// KeyName is the key of the field that holds the name to display.
	KeyName = "name"

	entrySeparator = ";"
	pairSeparator  = "="

	// valuePunctuation lists the characters that a value can hold besides
	// ASCII letters and digits. It covers email addresses and plain URLs.
	// The percent sign is left out on purpose: a percent-encoded value could
	// hide a separator or a control character from this parser and reveal it
	// to a client that decodes the value.
	valuePunctuation = "_-.~:/@+"

	// maxCombiningMarks is the number of combining marks that Sanitize keeps
	// in a row. The scripts that stack marks use a few of them at most.
	maxCombiningMarks = 4

	zeroWidthNonJoiner = '\u200c'
	zeroWidthJoiner    = '\u200d'
	brailleBlank       = '\u2800'
)

// Field is one key-value pair of a structured moniker.
type Field struct {
	Key   string
	Value string
}

// Moniker is a parsed moniker.
type Moniker struct {
	raw     string
	fields  []Field
	ignored int
}

// Parse reads a moniker. It never fails: a string that holds no valid
// key-value pair is a plain moniker. A string longer than MaxLength is not
// a valid moniker and no node accepts it from a peer: it is not parsed,
// and only its first MaxLength bytes are shown.
//
// The rules of the structured format are:
//   - The entries are separated by ';', and an empty entry is skipped.
//   - The key and the value of an entry are separated by the first '='.
//   - A key holds ASCII letters, digits and '_'. It is not case-sensitive.
//   - A value holds ASCII letters, digits and the characters "_-.~:/@+".
//   - An entry with an empty or invalid key or value is ignored.
//   - When a key is repeated, the first entry is kept and the others are ignored.
//   - No space is trimmed: a space makes the key or the value invalid.
func Parse(raw string) Moniker {
	mon := Moniker{raw: raw}

	if len(raw) > MaxLength {
		return mon
	}

	for entry := range strings.SplitSeq(raw, entrySeparator) {
		if entry == "" {
			continue
		}

		key, value, found := strings.Cut(entry, pairSeparator)
		if !found || !isValidKey(key) || !isValidValue(value) {
			mon.ignored++

			continue
		}

		key = lowerASCII(key)
		if _, exists := mon.Field(key); exists {
			mon.ignored++

			continue
		}

		mon.fields = append(mon.fields, Field{Key: key, Value: value})
	}

	return mon
}

// Raw returns the moniker as the node declared it.
func (m Moniker) Raw() string {
	return m.raw
}

// IsStructured reports whether the moniker holds at least one valid key-value pair.
func (m Moniker) IsStructured() bool {
	return len(m.fields) > 0
}

// Fields returns the valid key-value pairs in the order they were declared.
// The keys are in lower case.
func (m Moniker) Fields() []Field {
	fields := make([]Field, len(m.fields))
	copy(fields, m.fields)

	return fields
}

// Field returns the value of the given key. The key must be in lower case.
func (m Moniker) Field(key string) (string, bool) {
	for _, field := range m.fields {
		if field.Key == key {
			return field.Value, true
		}
	}

	return "", false
}

// Ignored returns the number of entries that are not empty and that the
// parser ignored. It is zero for a well-formed structured moniker.
func (m Moniker) Ignored() int {
	return m.ignored
}

// DisplayName returns the name to show for the node. It is the value of the
// "name" field, or the sanitized moniker if there is no such field.
// It is never longer than MaxLength bytes.
func (m Moniker) DisplayName() string {
	if name, ok := m.Field(KeyName); ok {
		return name
	}

	return Sanitize(truncate(m.raw))
}

// Sanitize makes a string that a node declared safe to show next to other text.
//
// It removes the characters that are not shown or that change how the
// surrounding text is shown: the control characters, the line and paragraph
// separators, the invisible format characters, including the bidirectional
// controls, and the characters that are shown as a blank, like the Hangul
// fillers. It also removes the invalid UTF-8 bytes.
// It replaces any kind of space with a plain space and trims the result.
// It keeps maxCombiningMarks combining marks in a row and removes the others,
// so a stack of marks can't be drawn over the text around.
//
// The zero width joiner and non-joiner are kept, because Persian and several
// other scripts need them to be written correctly. A string that holds nothing
// to show gives an empty string.
//
// The result can still hold characters that look like other characters:
// Sanitize does not make a name trustworthy.
func Sanitize(str string) string {
	var buf strings.Builder
	buf.Grow(len(str))

	marks := 0
	visible := false
	for _, char := range str {
		switch {
		case isJoiner(char):
			// A joiner is kept, and it doesn't start a new run of combining marks.

		case isUnsafeRune(char):
			// An invalid byte comes here as utf8.RuneError.
			continue

		case unicode.Is(unicode.Zs, char):
			char = ' '
			marks = 0

		case unicode.In(char, unicode.Mn, unicode.Me):
			marks++
			if marks > maxCombiningMarks {
				continue
			}

		default:
			marks = 0
			visible = true
		}

		buf.WriteRune(char)
	}

	if !visible {
		return ""
	}

	return strings.Trim(buf.String(), " ")
}

func isJoiner(r rune) bool {
	return r == zeroWidthNonJoiner || r == zeroWidthJoiner
}

func isUnsafeRune(r rune) bool {
	// A valid U+FFFD can't be told apart from an invalid byte: both are removed.
	if r == utf8.RuneError || r == brailleBlank {
		return true
	}

	return unicode.In(r, unicode.Cc, unicode.Cf, unicode.Zl, unicode.Zp,
		unicode.Other_Default_Ignorable_Code_Point)
}

// truncate cuts a string to MaxLength bytes, without cutting a character.
func truncate(str string) string {
	if len(str) <= MaxLength {
		return str
	}

	cut := MaxLength
	for cut > 0 && !utf8.RuneStart(str[cut]) {
		cut--
	}

	return str[:cut]
}

func isValidKey(key string) bool {
	if key == "" {
		return false
	}

	for i := range len(key) {
		if !isASCIIAlphanumeric(key[i]) && key[i] != '_' {
			return false
		}
	}

	return true
}

func isValidValue(value string) bool {
	if value == "" {
		return false
	}

	for i := range len(value) {
		if !isASCIIAlphanumeric(value[i]) && strings.IndexByte(valuePunctuation, value[i]) < 0 {
			return false
		}
	}

	return true
}

func isASCIIAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// lowerASCII converts the ASCII letters of a valid key to lower case.
// The Unicode case mapping is not used: it maps some non-ASCII characters to
// ASCII letters, like the Kelvin sign to 'k'.
func lowerASCII(key string) string {
	lower := []byte(key)
	for i, c := range lower {
		if c >= 'A' && c <= 'Z' {
			lower[i] = c + ('a' - 'A')
		}
	}

	return string(lower)
}
