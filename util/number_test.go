package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatNumberWithDelimiters(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{123456, "123,456"},
		{123, "123"},
		{0, "0"},
		{-123, "-123"},
		{-123456, "-123,456"},
	}

	for _, test := range tests {
		result := FormatIntWithDelimiters(test.input)
		assert.Equal(t, test.expected, result, "FormatNumber(%d)", test.input)
	}
}

func TestFormatFloatWithDelimiters(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{123456.789, "123,456.789"},
		{123.456, "123.456"},
		{0.0, "0"},
		{-123.456, "-123.456"},
		{-123456.789, "-123,456.789"},

		{123456, "123,456"},
		{123, "123"},
		{0, "0"},
		{-123, "-123"},
		{-123456, "-123,456"},
	}

	for _, test := range tests {
		result := FormatFloatWithDelimiters(test.input, -1)
		assert.Equal(t, test.expected, result, "FormatFloat(%f)", test.input)
	}
}
