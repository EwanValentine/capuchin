package main

import (
	"log"

	"github.com/EwanValentine/capuchin/query"
	"github.com/EwanValentine/capuchin/source"
)

func main() {
	s := source.NewFileSource()
	fileSource, err := s.Load("./query/test-data.csv")
	if err != nil {
		log.Panic(err)
	}

	query := &query.Query{
		Select: []string{"user_id", "date"},
		Where:  "user_id = abc123",
	}
	query.Source(fileSource)

	results, err := query.Exec()
	if err != nil {
		log.Panic(err)
	}

	log.Println(results)
}
