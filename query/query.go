package query

import (
	"encoding/csv"
	"io"
	"strings"
)

// Query -
type Query struct {
	// Select is a list of column names as strings
	Select []string `json:"select"`
	// Where is a SQL like where clause `col = abc-123`
	Where string `json:"where"`
	// Limit of rows returned
	Limit  *int `json:"limit"`
	reader *csv.Reader
}

// Exec query
func (q *Query) Exec() (Results, error) {
	count := 0
	selectCols := map[string]int{}
	results := Results{}

	record, err := q.reader.Read()
	if err != nil {
		return results, err
	}

	if err == io.EOF {
		return results, nil
	}

	for key, value := range record {
		for _, selector := range q.Select {
			if value == selector {
				selectCols[selector] = key
			}
		}
	}

	for {
		// @todo - this is probably shit, see: https://stackoverflow.com/questions/67685288/how-to-filter-csv-file-into-columns-on-go
		tmpResults := Results{}
		record, err := q.reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return results, err
		}

		// If select has rows defined
		if len(q.Select) > 0 {
			var row Row

			// For each value in the record, or csv row
			for key, value := range record {
				// For each column number in the selected columns
				for columnName, columnNumber := range selectCols {
					if key == columnNumber {
						row = append(row, Result{
							Key:   columnName,
							Value: value,
						})
					}
				}
			}

			// If the row count is filled, add it to the results
			if len(row) > 0 {
				tmpResults = append(tmpResults, row)
			}
		}

		// Where
		if q.Where != "" {
			tmpResults = q.where(results, tmpResults)
			results = tmpResults
		} else {
			results = append(results, tmpResults...)
		}

		count++

		if q.Limit != nil {
			if count >= *q.Limit {
				break
			}
		}
	}

	return results, nil
}

func (q *Query) where(results, tmpResults Results) Results {
	// Split where clause into column and value
	parts := strings.Split(q.Where, "=")

	columnName := strings.TrimSpace(parts[0])
	whereValue := strings.TrimSpace(parts[1])

	for _, result := range tmpResults {
		var filtered Row
		for _, row := range result {
			if row.Key == columnName && row.Value == whereValue {
				filtered = append(filtered, row)
			}
		}

		if len(filtered) > 0 {
			results = append(results, filtered)
		}
	}

	return results
}

// Source - csv reader
func (q *Query) Source(r *csv.Reader) *Query {
	q.reader = r
	return q
}

// Result -
type Results []Row

// Row -
type Row []Result

// Row -
type Result struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
