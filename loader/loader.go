// loader loads the data in raw csv format, a date key is provided
// the date key is used to create a range. The range is split by the number
// of nodes. Each split becomes a shard, which is represented by a single node.
//
// Each shard is loaded into the memory of each node. The queries are ran across
// all nodes, the results are combined and returned.
package loader

import (
	"encoding/csv"
	"log"

	"github.com/EwanValentine/capuchin/cluster"
)

// DataSource -
type DataSource interface {
	Read(start, end int) (*csv.Reader, error)
	Write(reader *csv.Reader) error
}

// DistributionManager -
type DistributionManager interface {
	NodeRemoved() <-chan cluster.Node
	NodeAdded() <-chan cluster.Node
}

// Loader -
type Loader struct {
	// DataSource is the service that actually grabs the shard of data from the raw data
	DataSource DataSource
	// DistributionManager is the integration with etcd
	DistributionManager DistributionManager
	// ShardID is 'which shard am I?'
	ShardID   int
	Start     int
	End       int
	nodeCount int
	stop      chan struct{}
}

// NewLoader -
func NewLoader(dataSource DataSource, distributionManager DistributionManager, shardID, nodeCount int) *Loader {
	return &Loader{
		DataSource:          dataSource,
		DistributionManager: distributionManager,
		nodeCount:           nodeCount,
		stop:                make(chan struct{}),
		ShardID:             shardID,
	}
}

// Load takes an input data stream, and writes that data into the raw data store.
func (l *Loader) Load(start, end int) error {
	l.DataSource.Read(start, end)
	return nil
}

// Write -
func (l *Loader) Write(reader *csv.Reader) error {
	return l.DataSource.Write(reader)
}

// AssignShardID -
func (l *Loader) AssignShardID(shardID int) {
	l.ShardID = shardID
}

// Shard takes the maximum start and end date of the entire raw data set,
// it generates a sequence of periods, gets the start and end date for just
// this nodes period. It then 'rebalances', which means it reloads the data
// using the new start and end date for just this node.
func (l *Loader) Shard(start, end int) {
	l.Start = start
	l.End = end
	go func() {
		for {
			select {
			case node := <-l.DistributionManager.NodeRemoved():
				log.Printf("Node %d deleted", node.ID)
				l.nodeCount--
				periods := GeneratePeriods(l.nodeCount, start, end)
				log.Println("New periods: ", periods)
				l.rebalanceDecreased(periods, node.ID)
			case node := <-l.DistributionManager.NodeAdded():
				log.Printf("Node %d added", node.ID)
				l.nodeCount++
				periods := GeneratePeriods(l.nodeCount, start, end)
				log.Println("New periods: ", periods)
				l.rebalanceIncreased(periods, node.ID)
			case <-l.stop:
				return
			default:
				continue
			}
		}
	}()
}

func (l *Loader) rebalanceDecreased(periods []Period, nodeID int) {
	shardID := l.ShardID

	// Decrement the shard ID if shard ID is greater than the
	// removed shard id
	if l.ShardID != 0 && l.ShardID > nodeID {
		shardID--
		l.AssignShardID(shardID)
	} else {
		l.AssignShardID(0)
	}

	// Get the periods for the new shard ID
	period := periods[l.ShardID]

	l.Start = period.Start
	l.End = period.End

	// Load the new period into memory
	l.Load(period.Start, period.End)
}

func (l *Loader) rebalanceIncreased(periods []Period, nodeID int) {
	shardID := l.ShardID + 1
	l.AssignShardID(shardID)

	// Get the periods for the new shard ID
	period := periods[shardID]

	l.Start = period.Start
	l.End = period.End

	// Load the new period into memory
	l.Load(period.Start, period.End)
}

// Stop -
func (l *Loader) Stop() {
	l.stop <- struct{}{}
}
