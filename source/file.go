package source

import (
	"encoding/csv"
	"os"
)

// FileSource -
type FileSource struct {
	prefix string
}

// NewFileSource - prefix example `./test-data/`
func NewFileSource(prefix string) *FileSource {
	return &FileSource{prefix: prefix}
}

// Read -
func (s *FileSource) Read(start, end int) (*csv.Reader, error) {
	f, err := os.Open(s.prefix + generateFilePath(start, end))
	reader := csv.NewReader(f)
	return reader, err
}

// Write -
func (s *FileSource) Write([]byte) error {
	return nil
}
