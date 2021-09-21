package query

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name     string
	Query    Query
	File     string
	Expected []Result
}

var (
	testCases = []testCase{{
		Name: "test where",
		Query: Query{
			Select: []string{"user_id"},
			Where:  "user_id = abc123",
		},
		File: "test-data.csv",
		Expected: []Result{{
			Key:   "user_id",
			Value: "abc123",
		}},
	}, {
		Name: "test can return all if no where",
		Query: Query{
			Select: []string{"order_id", "user_id", "date"},
		},
		File: "test-data.csv",
		Expected: []Result{{
			Key:   "order_id",
			Value: "abc123",
		}, {
			Key:   "user_id",
			Value: "abc123",
		}, {
			Key:   "date",
			Value: "2021-09-01",
		}},
	}}
)

func TestCases(t *testing.T) {
	for _, testCase := range testCases {
		f, err := os.Open("./" + testCase.File)
		if err != nil {
			log.Panic(err)
		}

		testCase.Query.Source(csv.NewReader(f))
		actual, err := testCase.Query.Exec()
		require.NoError(t, err)
		if !reflect.DeepEqual(actual, testCase.Expected) {
			t.Errorf("%s: expected %v, got %v", testCase.Name, testCase.Expected, actual)
		}
	}
}

func BenchmarkQueryWhere_10(b *testing.B) {
	testCase := testCases[0]
	f, err := os.Open("./" + testCase.File)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < b.N; i++ {
		testCase.Query.Source(csv.NewReader(f))
		testCase.Query.Exec()
	}
}

func BenchmarkQuerySelectAll_10(b *testing.B) {
	testCase := testCases[1]
	f, err := os.Open("./" + testCase.File)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < b.N; i++ {
		testCase.Query.Source(csv.NewReader(f))
		testCase.Query.Exec()
	}
}
