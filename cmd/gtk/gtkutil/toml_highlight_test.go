//go:build gtk

package gtkutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeTOMLHighlightRanges_commentLine(t *testing.T) {
	content := "# full-line comment\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 1)
	assert.Equal(t, TOMLTagComment, ranges[0].Tag)
	assert.Equal(t, 0, ranges[0].Start)
	assert.Equal(t, len("# full-line comment"), ranges[0].End)
}

func TestComputeTOMLHighlightRanges_indentedComment(t *testing.T) {
	content := "  # indented comment\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 1)
	assert.Equal(t, TOMLTagComment, ranges[0].Tag)
	assert.Equal(t, 2, ranges[0].Start)
	assert.Equal(t, len(content)-1, ranges[0].End)
}

func TestComputeTOMLHighlightRanges_tableHeader(t *testing.T) {
	content := "[foo]\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 1)
	assert.Equal(t, TOMLTagTable, ranges[0].Tag)
	assert.Equal(t, 0, ranges[0].Start)
	assert.Equal(t, len(content)-1, ranges[0].End)
}

func TestComputeTOMLHighlightRanges_indentedTableHeader(t *testing.T) {
	content := "  [foo]\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 1)
	assert.Equal(t, TOMLTagTable, ranges[0].Tag)
	assert.Equal(t, 2, ranges[0].Start)
	assert.Equal(t, len(content)-1, ranges[0].End)
}

func TestComputeTOMLHighlightRanges_keyValue(t *testing.T) {
	content := "key = 64\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 2)
	assert.Equal(t, TOMLHighlightRange{0, len("key"), TOMLTagKey}, ranges[0])
	assert.Equal(t, TOMLHighlightRange{
		len("key = "),
		len("key = 64"),
		TOMLTagNumber,
	}, ranges[1])
}

func TestComputeTOMLHighlightRanges_indentedKeyValue(t *testing.T) {
	content := "  key = 'value'\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 2)
	assert.Equal(t, TOMLHighlightRange{2, 2 + len("key"), TOMLTagKey}, ranges[0])
	assert.Equal(t, TOMLHighlightRange{
		2 + len("key = "),
		len(content) - 1,
		TOMLTagString,
	}, ranges[1])
}

func TestComputeTOMLHighlightRanges_boolean(t *testing.T) {
	content := "key = false\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 2)
	assert.Equal(t, TOMLTagKey, ranges[0].Tag)
	assert.Equal(t, TOMLTagBoolean, ranges[1].Tag)
	assert.Equal(t, len("key = "), ranges[1].Start)
	assert.Equal(t, len("key = false"), ranges[1].End)
}

func TestComputeTOMLHighlightRanges_inlineComment(t *testing.T) {
	content := "key = 10  # comments\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.GreaterOrEqual(t, len(ranges), 3)

	var comment *TOMLHighlightRange
	for i := range ranges {
		if ranges[i].Tag == TOMLTagComment {
			comment = &ranges[i]

			break
		}
	}

	require.NotNil(t, comment)
	assert.Equal(t, len("key = 10  "), comment.Start)
	assert.Equal(t, len(content)-1, comment.End)
}

func TestComputeTOMLHighlightRanges_multilineOffsets(t *testing.T) {
	content := "[foo]\nkey = 64\n"
	ranges := ComputeTOMLHighlightRanges(content)

	require.Len(t, ranges, 3)
	assert.Equal(t, TOMLTagTable, ranges[0].Tag)
	assert.Equal(t, 0, ranges[0].Start)

	line2Start := len("[foo]\n")
	assert.Equal(t, TOMLTagKey, ranges[1].Tag)
	assert.Equal(t, line2Start, ranges[1].Start)
	assert.Equal(t, line2Start+len("key"), ranges[1].End)

	assert.Equal(t, TOMLTagNumber, ranges[2].Tag)
	assert.Equal(t, line2Start+len("key = "), ranges[2].Start)
}

func TestFindTOMLCommentIndex_ignoresHashInStrings(t *testing.T) {
	assert.Equal(t, -1, findTOMLCommentIndex(`key = "#not-a-comment"`))
	assert.Equal(t, len(`key = "" `), findTOMLCommentIndex(`key = "" # real comment`))
}

func TestIsTOMLNumber(t *testing.T) {
	assert.True(t, isTOMLNumber("64"))
	assert.True(t, isTOMLNumber("-1"))
	assert.True(t, isTOMLNumber("3.14"))
	assert.False(t, isTOMLNumber(""))
	assert.False(t, isTOMLNumber("abc"))
	assert.False(t, isTOMLNumber("true"))
}
