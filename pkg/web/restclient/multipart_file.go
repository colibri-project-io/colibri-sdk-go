package restclient

import "io"

type MultipartFile struct {
	FileName    string
	File        io.Reader
	ContentType string
}
