package loader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanGenerateNewShardRange(t *testing.T) {
	shardCount := 4
	startDate := 20190101
	endDate := 20210910
	periods := GeneratePeriods(shardCount, startDate, endDate)

	require.Len(t, periods, shardCount)
	require.Equal(t, periods[0].Start, startDate)
	require.Equal(t, periods[3].End, endDate)
}
