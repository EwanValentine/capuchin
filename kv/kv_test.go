package kv

import (
	"log"
	"testing"

	"github.com/EwanValentine/capuchin/source"
	"github.com/stretchr/testify/require"
)

func TestCanLoadCSV(t *testing.T) {
	store, err := New()
	require.NoError(t, err)

	start := 20190101
	end := 20190201

	fileSource := source.NewFileSource("../test-data/")
	reader, err := fileSource.Read(start, end)
	require.NoError(t, err)

	err = store.LoadCSV(reader)
	require.NoError(t, err)

	result, err := store.Get("20190101")
	require.NoError(t, err)

	log.Println(string(result))
}
