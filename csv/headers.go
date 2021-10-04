package csv

import (
	"encoding/csv"
	"io"
)

// Columns is a reference to each header columns
// name, and their corresponding index number,
// for example: "date": 0.
type Columns map[string]int

// ExtractHeaders - takes the first row of a csv reader, returns that row as a key value pair
// with a reference to the column index number, and returns the original reader again
func ExtractHeaders(reader *csv.Reader) (*csv.Reader, Columns, error) {
	cols := make(Columns, 0)

	// Read the first row from the csv reader
	record, err := reader.Read()
	if err != nil {
		return reader, cols, err
	}

	// If EOF is reached, return the reader and an empty map
	if err == io.EOF {
		return reader, cols, nil
	}

	// Iterate through each value in the record
	for key, value := range record {
		cols[value] = key
	}

	return reader, cols, nil
}
