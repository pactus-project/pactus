package downloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChunks(t *testing.T) {
	tests := []struct {
		contentLength int64
		totalChunks   int64
		expected      []*chunk
	}{
		{
			contentLength: 181403648,
			totalChunks:   16,
			expected: []*chunk{
				{start: 0, end: 11337727},
				{start: 11337728, end: 22675455},
				{start: 22675456, end: 34013183},
				{start: 34013184, end: 45350911},
				{start: 45350912, end: 56688639},
				{start: 56688640, end: 68026367},
				{start: 68026368, end: 79364095},
				{start: 79364096, end: 90701823},
				{start: 90701824, end: 102039551},
				{start: 102039552, end: 113377279},
				{start: 113377280, end: 124715007},
				{start: 124715008, end: 136052735},
				{start: 136052736, end: 147390463},
				{start: 147390464, end: 158728191},
				{start: 158728192, end: 170065919},
				{start: 170065920, end: 181403647},
			},
		},
		{
			contentLength: 10,
			totalChunks:   3,
			expected: []*chunk{
				{start: 0, end: 2},
				{start: 3, end: 5},
				{start: 6, end: 9},
			},
		},
		{
			contentLength: 10,
			totalChunks:   1,
			expected: []*chunk{
				{start: 0, end: 9},
			},
		},
		{
			contentLength: 0,
			totalChunks:   1,
			expected: []*chunk{
				{start: 0, end: -1},
			},
		},
	}

	for _, tt := range tests {
		actual := createChunks(tt.contentLength, tt.totalChunks)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestChunkRangeHeader(t *testing.T) {
	tests := []struct {
		chunk    chunk
		expected string
	}{
		{
			chunk:    chunk{start: 0, end: 499},
			expected: "bytes=0-499",
		},
		{
			chunk:    chunk{start: 500, end: 999},
			expected: "bytes=500-999",
		},
	}

	for _, tt := range tests {
		actual := tt.chunk.rangeHeader()
		assert.Equal(t, tt.expected, actual)
	}
}

func TestChunkSize(t *testing.T) {
	tests := []struct {
		chunk    chunk
		expected int64
	}{
		{
			chunk:    chunk{start: 0, end: 499},
			expected: 500,
		},
		{
			chunk:    chunk{start: 500, end: 999},
			expected: 500,
		},
		{
			chunk:    chunk{start: 0, end: 0},
			expected: 1,
		},
	}

	for _, tt := range tests {
		actual := tt.chunk.size()
		assert.Equal(t, tt.expected, actual)
	}
}
