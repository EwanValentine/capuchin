// Ingests CSV data, with a given date key to use as the primary key.
// Data is then stored in a 'location', which is what Capuchin uses as
// the data source.
package ingest

// New -
func New(location string, dateKey string, adapter Adapter) *Adapter {
	return &Ingest{
		location: location,
		dateKey:  dateKey,
	}
}

// Task
type Task struct {
	location string
	dateKey  string
}

// Ingest -
func (i *Task) Ingest(location string, dateKey string, csvData []byte) error {
	return nil
}
