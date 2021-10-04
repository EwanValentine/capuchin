package ingest

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	csvutils "github.com/EwanValentine/capuchin/csv"
	"github.com/EwanValentine/capuchin/dates"
)

// FileSource -
type FileSource struct {
	path        string
	dateKey     string
	dateFormat  string
	fileHandler *os.File
	csvReader   *csv.Reader
}

// NewFileSource -
func NewFileSource(path, dateKey, dateFormat string) *FileSource {
	return &FileSource{
		path:        path,
		dateKey:     dateKey,
		dateFormat:  dateFormat,
		fileHandler: nil,
		csvReader:   nil,
	}
}

// Read -
func (f *FileSource) Read(path string) (*csv.Reader, error) {
	fileHandler, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file with path %s: %v", path, err)
	}

	f.fileHandler = fileHandler

	reader := csv.NewReader(f.fileHandler)

	return reader, nil
}

// Close file handler
func (f *FileSource) Close() error {
	return f.fileHandler.Close()
}

// Read reads
func (f *FileSource) Write(reader *csv.Reader) error {
	handler, err := os.OpenFile(f.path, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	// Pass file handler to new CSV writer
	csvWriter := csv.NewWriter(handler)
	if err != nil {
		return err
	}

	// Extract the header values, in format "date": 0, etc
	reader, headers, err := csvutils.ExtractHeaders(reader)
	if err != nil {
		return fmt.Errorf("error extracting headers: %v", err)
	}

	// Get the string values for each header, so "date", "user_id" etc...
	headerValues := make([]string, len(headers))
	for k, v := range headers {
		headerValues[v] = k
	}

	// Write the headers as the first line of the output file
	if err := csvWriter.Write(headerValues); err != nil {
		return fmt.Errorf("error writing headers to CSV file: %v", err)
	}

	// Get the date's index
	dateIdx, ok := headers[f.dateKey]
	if !ok {
		return fmt.Errorf("error finding date index with name %s: %v", f.dateKey, err)
	}

	// Iterate over each row of the CSV input reader
	for {
		var tmpRow []string

		// Read each line from the CSV reader
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Get the value for the date, using the index
		date := record[dateIdx]

		// Parse the date
		d, err := dates.ParseDate(date, f.dateFormat)
		if err != nil {
			return fmt.Errorf("error parsing date %s: %v", date, err)
		}

		// Date always gets placed into records first, as index 0
		tmpRow = append(tmpRow, strconv.Itoa(d))

		// Rest of the fields are then added
		tmpRow = append(tmpRow, record[1:]...)

		// Write the new formatted row into the writer
		if err := csvWriter.Write(tmpRow); err != nil {
			return fmt.Errorf("error writing row to new csv writer: %v", err)
		}
	}

	// Flush to file handler
	csvWriter.Flush()
	f.Close()

	return nil
}
