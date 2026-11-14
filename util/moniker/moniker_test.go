package moniker

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		fields      []Field
		ignored     int
		displayName string
	}{
		{
			name:        "Empty moniker",
			raw:         "",
			displayName: "",
		},
		{
			name:        "Plain moniker",
			raw:         "John Doe",
			ignored:     1,
			displayName: "John Doe",
		},
		{
			name:        "Plain moniker with spaces",
			raw:         "My Pactus Node",
			ignored:     1,
			displayName: "My Pactus Node",
		},
		{
			name: "Structured moniker of PIP-45",
			raw:  "name=john_doe;email=john.doe@example.com;web=https://example.com",
			fields: []Field{
				{Key: "name", Value: "john_doe"},
				{Key: "email", Value: "john.doe@example.com"},
				{Key: "web", Value: "https://example.com"},
			},
			displayName: "john_doe",
		},
		{
			name:        "Name is not the first field",
			raw:         "web=https://example.com;name=alice",
			fields:      []Field{{Key: "web", Value: "https://example.com"}, {Key: "name", Value: "alice"}},
			displayName: "alice",
		},
		{
			name:        "No name field",
			raw:         "web=https://example.com",
			fields:      []Field{{Key: "web", Value: "https://example.com"}},
			displayName: "web=https://example.com",
		},
		{
			name:        "Keys are not case-sensitive",
			raw:         "NAME=Alice;Web_Site=example.com",
			fields:      []Field{{Key: "name", Value: "Alice"}, {Key: "web_site", Value: "example.com"}},
			displayName: "Alice",
		},
		{
			name:        "First entry wins for a repeated key",
			raw:         "name=alice;name=bob;Name=carol",
			fields:      []Field{{Key: "name", Value: "alice"}},
			ignored:     2,
			displayName: "alice",
		},
		{
			name:        "Empty entries are skipped",
			raw:         ";;name=alice;;web=example.com;",
			fields:      []Field{{Key: "name", Value: "alice"}, {Key: "web", Value: "example.com"}},
			displayName: "alice",
		},
		{
			name:        "Entry without a separator is ignored",
			raw:         "name=alice;hello",
			fields:      []Field{{Key: "name", Value: "alice"}},
			ignored:     1,
			displayName: "alice",
		},
		{
			name:        "Empty key is ignored",
			raw:         "=alice;web=example.com",
			fields:      []Field{{Key: "web", Value: "example.com"}},
			ignored:     1,
			displayName: "=alice;web=example.com",
		},
		{
			name:        "Empty value is ignored",
			raw:         "name=;web=example.com",
			fields:      []Field{{Key: "web", Value: "example.com"}},
			ignored:     1,
			displayName: "name=;web=example.com",
		},
		{
			name:        "Second separator makes the value invalid",
			raw:         "name=alice=bob;web=example.com",
			fields:      []Field{{Key: "web", Value: "example.com"}},
			ignored:     1,
			displayName: "name=alice=bob;web=example.com",
		},
		{
			name:        "Spaces are not trimmed",
			raw:         "web=example.com; name=alice;email =a@b.c;key= value",
			fields:      []Field{{Key: "web", Value: "example.com"}},
			ignored:     3,
			displayName: "web=example.com; name=alice;email =a@b.c;key= value",
		},
		{
			name:        "Space in the value is invalid",
			raw:         "name=My Node",
			ignored:     1,
			displayName: "name=My Node",
		},
		{
			name:        "Invalid key characters",
			raw:         "na-me=alice;na.me=bob;nam\u00e9=carol",
			ignored:     3,
			displayName: "na-me=alice;na.me=bob;nam\u00e9=carol",
		},
		{
			name:        "Kelvin sign is not the letter k",
			raw:         "\u212aey=value",
			ignored:     1,
			displayName: "\u212aey=value",
		},
		{
			name:        "Percent-encoded value is invalid",
			raw:         "name=alice%3Bname%3Dbob",
			ignored:     1,
			displayName: "name=alice%3Bname%3Dbob",
		},
		{
			name:        "Markup in the value is invalid",
			raw:         "name=<script>alert(1)</script>",
			ignored:     1,
			displayName: "name=<script>alert(1)</script>",
		},
		{
			name:        "Query string in the value is invalid",
			raw:         "web=https://example.com/?a=b&c=d",
			ignored:     1,
			displayName: "web=https://example.com/?a=b&c=d",
		},
		{
			name:        "Non-ASCII value is invalid",
			raw:         "name=\u062c\u0648\u0627\u062f",
			ignored:     1,
			displayName: "name=\u062c\u0648\u0627\u062f",
		},
		{
			name:        "Control characters are removed from a plain moniker",
			raw:         "\x1b[31mAlice\r\n\x00",
			ignored:     1,
			displayName: "[31mAlice",
		},
		{
			name:        "Control characters make a value invalid",
			raw:         "name=alice\n;web=example.com",
			fields:      []Field{{Key: "web", Value: "example.com"}},
			ignored:     1,
			displayName: "name=alice;web=example.com",
		},
		{
			name:        "All value punctuation",
			raw:         "key_1=aZ09_-.~:/@+",
			fields:      []Field{{Key: "key_1", Value: "aZ09_-.~:/@+"}},
			displayName: "key_1=aZ09_-.~:/@+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mon := Parse(tt.raw)

			assert.Equal(t, tt.raw, mon.Raw())
			assert.Equal(t, tt.ignored, mon.Ignored())
			assert.Equal(t, tt.displayName, mon.DisplayName())
			assert.Equal(t, len(tt.fields) > 0, mon.IsStructured())
			if len(tt.fields) == 0 {
				assert.Empty(t, mon.Fields())
			} else {
				assert.Equal(t, tt.fields, mon.Fields())
			}
		})
	}
}

func TestTooLongMoniker(t *testing.T) {
	t.Run("Longest moniker is parsed", func(t *testing.T) {
		raw := "name=" + strings.Repeat("a", MaxLength-5)
		mon := Parse(raw)

		assert.True(t, mon.IsStructured())
		assert.Equal(t, strings.Repeat("a", MaxLength-5), mon.DisplayName())
	})

	t.Run("Longer moniker is not parsed", func(t *testing.T) {
		raw := "name=alice;" + strings.Repeat("a", MaxLength)
		mon := Parse(raw)

		assert.Equal(t, raw, mon.Raw())
		assert.False(t, mon.IsStructured())
		assert.Zero(t, mon.Ignored())
		assert.Equal(t, raw[:MaxLength], mon.DisplayName())
	})

	t.Run("Display name is not cut inside a character", func(t *testing.T) {
		// 127 bytes, then a character of two bytes.
		raw := strings.Repeat("a", MaxLength-1) + "\u00e9\u00e9"
		mon := Parse(raw)

		assert.Equal(t, strings.Repeat("a", MaxLength-1), mon.DisplayName())
	})

	t.Run("Huge moniker", func(t *testing.T) {
		raw := strings.Repeat("k=v;", 1_000_000)
		mon := Parse(raw)

		assert.False(t, mon.IsStructured())
		assert.Len(t, mon.DisplayName(), MaxLength)
	})
}

func TestField(t *testing.T) {
	mon := Parse("Name=alice;web=example.com")

	name, ok := mon.Field("name")
	assert.True(t, ok)
	assert.Equal(t, "alice", name)

	_, ok = mon.Field("email")
	assert.False(t, ok)

	// The key of a parsed field is in lower case.
	_, ok = mon.Field("Name")
	assert.False(t, ok)
}

func TestFieldsIsACopy(t *testing.T) {
	mon := Parse("name=alice")

	fields := mon.Fields()
	fields[0].Value = "bob"

	assert.Equal(t, "alice", mon.DisplayName())
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Empty", input: "", expected: ""},
		{name: "ASCII", input: "Alice's node #1", expected: "Alice's node #1"},
		{name: "Other scripts", input: "N\u0153ud \u00c9lise \u8282\u70b9", expected: "N\u0153ud \u00c9lise \u8282\u70b9"},
		{name: "Control characters", input: "a\x00b\tc\nd\re\x1bf\x7fg\u0085h", expected: "abcdefgh"},
		{name: "Line and paragraph separators", input: "a\u2028b\u2029c", expected: "abc"},
		{name: "Bidirectional controls", input: "a\u202eb\u202dc\u2066d\u2069e\u200ef\u200fg\u061ch", expected: "abcdefgh"},
		{name: "Invisible format characters", input: "a\u200bb\u2060c\ufeffd\u00ade", expected: "abcde"},
		{name: "Tag characters", input: "a\U000e0041b", expected: "ab"},
		{
			name:     "Joiners are kept",
			input:    "\u0645\u06cc\u200c\u062e\u0648\u0627\u0647\u0645 \U0001f468\u200d\U0001f4bb",
			expected: "\u0645\u06cc\u200c\u062e\u0648\u0627\u0647\u0645 \U0001f468\u200d\U0001f4bb",
		},
		{name: "Invalid UTF-8", input: "a\xffb\xc3", expected: "ab"},
		{name: "Hangul fillers", input: "a\u3164b\u115fc\u1160d\uffa0e", expected: "abcde"},
		{name: "Other ignorable characters", input: "a\u034fb\u17b4c\u17b5d\u2065e", expected: "abcde"},
		{name: "Braille blank", input: "a\u2800b", expected: "ab"},
		{name: "Any space becomes a plain space", input: "a\u00a0b\u3000c\u2003d e", expected: "a b c d e"},
		{name: "Spaces are trimmed", input: " \u00a0 alice \u3000 ", expected: "alice"},
		{name: "Only spaces", input: " \u00a0\u3000 ", expected: ""},
		{name: "Only invisible characters", input: "\u3164\u2800\u200b\u202e", expected: ""},
		{name: "Only joiners", input: "\u200d\u200c \u200d", expected: ""},
		{name: "Only combining marks", input: "\u0301\u0301", expected: ""},
		{name: "Emoji variation selector is kept", input: "\u2764\ufe0f", expected: "\u2764\ufe0f"},
		{name: "Few combining marks are kept", input: "a\u0301\u0302e\u0303", expected: "a\u0301\u0302e\u0303"},
		{
			name:     "Stack of combining marks",
			input:    "a" + strings.Repeat("\u0301", 60) + "b\u0301",
			expected: "a" + strings.Repeat("\u0301", 4) + "b\u0301",
		},
		{
			name:     "A joiner doesn't start a new stack",
			input:    "a" + strings.Repeat("\u0301\u0301\u200d", 10),
			expected: "a\u0301\u0301\u200d\u0301\u0301" + strings.Repeat("\u200d", 9),
		},
		{
			name:     "A removed character doesn't start a new stack",
			input:    "a" + strings.Repeat("\u0301\u0301\u200b", 10),
			expected: "a" + strings.Repeat("\u0301", 4),
		},
		{name: "Replacement character", input: "a\ufffdb", expected: "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Sanitize(tt.input))
		})
	}
}

func FuzzParse(f *testing.F) {
	f.Add("")
	f.Add("John Doe")
	f.Add("name=john_doe;email=john.doe@example.com;web=https://example.com")
	f.Add("name=alice;name=bob;;=;a=;=b")
	f.Add("\x1b[31m\u202ename=alice\xff")
	f.Add(" \u00a0a\u0301\u0301\u0301\u0301\u0301\u200d\u3164\u2800 ")
	f.Add(strings.Repeat("k=v;", 64))

	f.Fuzz(func(t *testing.T, raw string) {
		mon := Parse(raw)

		seen := make(map[string]bool)
		for _, field := range mon.Fields() {
			assert.True(t, isValidKey(field.Key))
			assert.True(t, isValidValue(field.Value))
			assert.Equal(t, strings.ToLower(field.Key), field.Key)
			assert.False(t, seen[field.Key], "repeated key %q", field.Key)
			seen[field.Key] = true

			// Each field comes from one entry of the moniker.
			assert.LessOrEqual(t, len(field.Key)+len(field.Value)+1, len(raw))
		}

		name := mon.DisplayName()
		assert.True(t, utf8.ValidString(name))
		assert.LessOrEqual(t, len(name), len(raw))
		assert.LessOrEqual(t, len(name), MaxLength)
		assert.Equal(t, name, Sanitize(name), "the display name is not clean")
		assert.Equal(t, name, strings.TrimSpace(name), "the display name is not trimmed")
		assert.NotContains(t, name, strings.Repeat("\u0301", maxCombiningMarks+1))
		for _, char := range name {
			assert.False(t, isUnsafeRune(char) && !isJoiner(char), "unsafe character %U", char)
		}

		if len(raw) > MaxLength {
			assert.False(t, mon.IsStructured())
		}

		// The fields of a moniker rebuild a moniker with the same fields.
		entries := make([]string, 0, len(mon.Fields()))
		for _, field := range mon.Fields() {
			entries = append(entries, field.Key+"="+field.Value)
		}
		rebuilt := Parse(strings.Join(entries, ";"))
		assert.Equal(t, mon.Fields(), rebuilt.Fields())
		assert.Zero(t, rebuilt.Ignored())
	})
}
