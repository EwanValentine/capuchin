package query

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
)

// Query -
type Query struct {
	// Select is a list of column names as strings
	Select []string `json:"select"`
	// Where is a SQL like where clause `col = abc-123`
	Where  string `json:"where"`
	reader *csv.Reader
}

// Exec query
func (q *Query) Exec() ([]Result, error) {
	count := 0
	selectCols := map[string]int{}
	results := []Result{}
	for {
		tmpResults := []Result{}
		record, err := q.reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panic(err)
		}

		if count == 0 {
			for key, value := range record {
				for _, selector := range q.Select {
					if value == selector {
						selectCols[selector] = key
					}
				}
			}
			count++
			continue
		}

		// Filter the columns
		for key, value := range record {
			for k, sel := range selectCols {
				if key == sel {
					tmpResults = append(tmpResults, Result{
						Key:   k,
						Value: value,
					})
				}
			}
		}

		// Where
		if q.Where != "" {
			parts := strings.Split(q.Where, " = ")

			col := parts[0]
			v := parts[1]

			for _, result := range tmpResults {
				if result.Key == col && result.Value == v {
					results = append(results, result)
				}
			}
		} else {
			results = tmpResults
		}

		count++
	}

	return results, nil
}

// Source - csv reader
func (q *Query) Source(r *csv.Reader) *Query {
	q.reader = r
	return q
}

// Result -
type Result struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
