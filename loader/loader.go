// loader loads the data in raw csv format, a date key is provided
// the date key is used to create a range. The range is split by the number
// of nodes. Each split becomes a shard, which is represented by a single node.
//
// Each shard is loaded into the memory of each node. The queries are ran across
// all nodes, the results are combined and returned.
package loader

import "log"

// DataSource -
type DataSource interface {
	Read() ([]byte, error)
	Write([]byte) error
}

// Loader -
type Loader struct {
	DataSource DataSource
	dateKey    string
	nodeCount  int
}

// NewLoader -
func NewLoader(dataSource DataSource, dateKey string, nodeCount int) *Loader {
	return &Loader{dateKey: dateKey, DataSource: dataSource, nodeCount: nodeCount}
}

// Load takes an input data stream, and writes that data into the raw data store.
func (l *Loader) Load() error {

	return nil
}

// Shard splits the data by the date key, by the number of nodes in the network
// and stores each shard in the data store. Each node then receives a message
// to sync with the generated shards.
func (l *Loader) Shard() error {

	// Mock start and end
	start := 20190101
	end := 20210920

	// Dummy shard count
	shards := (start + end) / l.nodeCount

	log.Println(shards)

	// Each shard interval could be stored in something like Redis
	// each node could be connected by etcd or some shit
	// The shard interval or updates could be signalled to each node

	return nil
}
