package loader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	csvutils "github.com/EwanValentine/capuchin/csv"
	"github.com/EwanValentine/capuchin/dates"
)

// FileLoader -
type FileLoader struct {
	path       string
	dateFormat string
	dateKey    string
}

// NewFileLoader -
func NewFileLoader(path, dateFormat, dateKey string) *FileLoader {
	return &FileLoader{path, dateFormat, dateKey}
}

// Write logic needs to take unstructured CSV data, split the data into
// chunks, by the given periods, and write the data to the given path.
func (l *FileLoader) Write(p *csv.Reader) error {
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	// Pass file handler to new CSV writer
	csvWriter := csv.NewWriter(f)
	if err != nil {
		return err
	}

	// Extract the header values, in format "date": 0, etc
	reader, headers, err := csvutils.ExtractHeaders(p)
	if err != nil {
		return fmt.Errorf("error extracting headers: %v", err)
	}

	headerValues := make([]string, len(headers))
	for k, v := range headers {
		headerValues[v] = k
	}
	if err := csvWriter.Write(headerValues); err != nil {
		return fmt.Errorf("error writing headers to CSV file: %v", err)
	}

	// Get the date's index
	dateIdx, ok := headers[l.dateKey]
	if !ok {
		return fmt.Errorf("error finding date index with name %s: %v", l.dateKey, err)
	}

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
		d, err := dates.ParseDate(date, l.dateFormat)
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
