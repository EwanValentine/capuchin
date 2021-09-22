package loader

// Period represents the start and end dates for each shard
type Period struct {
	Shard int
	Start int
	End   int
}

// GeneratePeriods -
func GeneratePeriods(shardCount, startDate, endDate int) []Period {
	diff := (endDate - startDate) / shardCount

	periods := []Period{}
	shard := 0
	for i := startDate; i+diff <= endDate; i = i + diff {
		periodStart := i
		periodEnd := i + diff
		periods = append(periods, Period{shard, periodStart, periodEnd})
		shard++
	}

	// Ensure the last shard has the endDate, this accounts for the remainder
	periods[shardCount-1].End = endDate

	return periods
}
