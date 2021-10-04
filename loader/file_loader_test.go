package loader

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanWriteToFile(t *testing.T) {
	f, err := os.Open("./test-data/raw-data.csv")
	require.NoError(t, err)

	format := "2006-02-01"
	reader := csv.NewReader(f)
	l := NewFileLoader("./test-data/output-data.csv", format, "date")
	err = l.Write(reader)
	require.NoError(t, err)
}
