package downloader

import "errors"

var (
	ErrHeaderRequest  = errors.New("request header error")
	ErrSHA256Mismatch = errors.New("sha256 mismatch")
	ErrCreateDir      = errors.New("create dir error")
	ErrDoRequest      = errors.New("error doing request")
	ErrFileWriting    = errors.New("error writing file")
	ErrNewRequest     = errors.New("error creating request")
	ErrOpenFileExists = errors.New("error opening existing file")
)
