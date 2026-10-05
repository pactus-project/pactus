package account

import "errors"

var (
	// ErrInvalidAnchorFlag is returned when the anchor flag is not 1.
	ErrInvalidAnchorFlag = errors.New("invalid anchor flag")

	// ErrAnchorHashLength is returned when the anchor hash length is out of range.
	ErrAnchorHashLength = errors.New("anchor hash length out of range")

	// ErrAnchorURILength is returned when the anchor URI is longer than 128 bytes.
	ErrAnchorURILength = errors.New("anchor uri length out of range")

	// ErrInvalidAnchorURI is returned when the anchor URI is not valid UTF-8.
	ErrInvalidAnchorURI = errors.New("invalid anchor uri")

	// ErrInvalidAnchorDeposit is returned when the anchor deposit is out of range.
	ErrInvalidAnchorDeposit = errors.New("invalid anchor deposit")

	// ErrAnchorTrailingBytes is returned when bytes follow a complete anchor.
	ErrAnchorTrailingBytes = errors.New("trailing bytes after anchor")
)
