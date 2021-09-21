package source

import (
	"encoding/csv"
	"os"
)

// FileSource -
type FileSource struct{}

// NewFileSource -
func NewFileSource() *FileSource {
	return &FileSource{}
}

// Load -
func (s *FileSource) Load(src string) (*csv.Reader, error) {
	f, err := os.Open(src)
	reader := csv.NewReader(f)
	return reader, err
}
