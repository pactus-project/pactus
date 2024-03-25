// This file contains code modified from the btcd project,
// which is licensed under the ISC License.
//
// Original license: https://github.com/btcsuite/btcd/blob/master/LICENSE
//

package bech32m

import (
	"fmt"

	"github.com/pactus-project/pactus/util/errors"
)

// MixedCaseError is returned when the bech32 string has both lower and uppercase
// characters.
type MixedCaseError struct{}

func (e MixedCaseError) Error() string {
	return "string not all lowercase or all uppercase"
}

func (e MixedCaseError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidBitGroupsError is returned when conversion is attempted between byte
// slices using bit-per-element of unsupported value.
type InvalidBitGroupsError struct{}

func (e InvalidBitGroupsError) Error() string {
	return "only bit groups between 1 and 8 allowed"
}

func (e InvalidBitGroupsError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidIncompleteGroupError is returned when then byte slice used as input has
// data of wrong length.
type InvalidIncompleteGroupError struct{}

func (e InvalidIncompleteGroupError) Error() string {
	return "invalid incomplete group"
}

func (e InvalidIncompleteGroupError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidLengthError is returned when the bech32 string has an invalid length
// given the BIP-173 defined restrictions.
type InvalidLengthError int

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("invalid bech32 string length %d", int(e))
}

func (e InvalidLengthError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidCharacterError is returned when the bech32 string has a character
// outside the range of the supported charset.
type InvalidCharacterError rune

func (e InvalidCharacterError) Error() string {
	return fmt.Sprintf("invalid character in string: '%c'", rune(e))
}

func (e InvalidCharacterError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidSeparatorIndexError is returned when the separator character '1' is
// in an invalid position in the bech32 string.
type InvalidSeparatorIndexError int

func (e InvalidSeparatorIndexError) Error() string {
	return fmt.Sprintf("invalid separator index %d", int(e))
}

func (e InvalidSeparatorIndexError) Code() int {
	return errors.ErrInvalidAddress
}

// NonCharsetCharError is returned when a character outside the specific
// bech32 charset is used in the string.
type NonCharsetCharError rune

func (e NonCharsetCharError) Error() string {
	return fmt.Sprintf("invalid character not part of charset: %v", int(e))
}

func (e NonCharsetCharError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidChecksumError is returned when the extracted checksum of the string
// is different than what was expected.
type InvalidChecksumError struct {
	Expected string
	Actual   string
}

func (e InvalidChecksumError) Error() string {
	return fmt.Sprintf("invalid checksum (expected %v got %v)",
		e.Expected, e.Actual)
}

func (e InvalidChecksumError) Code() int {
	return errors.ErrInvalidAddress
}

// InvalidDataByteError is returned when a byte outside the range required for
// conversion into a string was found.
type InvalidDataByteError byte

func (e InvalidDataByteError) Error() string {
	return fmt.Sprintf("invalid data byte: %v", byte(e))
}

func (e InvalidDataByteError) Code() int {
	return errors.ErrInvalidAddress
}
