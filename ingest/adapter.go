package ingest

// Adapter -
type Adapter interface {
	Load(path string, dateKey string) ([]byte, error)
}
