// loader loads the data in raw csv format, a date key is provided
// the date key is used to create a range. The range is split by the number
// of nodes. Each split becomes a shard, which is represented by a single node.
//
// Each shard is loaded into the memory of each node. The queries are ran across
// all nodes, the results are combined and returned.
package loader

import (
	"log"

	"github.com/EwanValentine/capuchin/cluster"
)

// DataSource -
type DataSource interface {
	Read(start, end int) ([]byte, error)
	Write([]byte) error
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

// AssignShardID -
func (l *Loader) AssignShardID(shardID int) {
	l.ShardID = shardID
}

// Shard splits the data by the date key, by the number of nodes in the network
// and stores each shard in the data store. Each node then receives a message
// to sync with the generated shards.
//
// Each shard interval could be stored in something like Redis
// each node could be connected by etcd or some shit
// The shard interval or updates could be signalled to each node

// Should have a shard size, say if the raw data is 1m rows,
// each shard should contain N number of rows. So if the shard size is
// 100,000 then we need 10 nodes. Then the start and end dates from each
// is stored along with the node's metadata in the service discovery stuff.
//
// So like... /nodes/1/20180101/20210910 - the shard size could just be the number
// of nodes in the cluster. So if you deploy 10 nodes, then each will have
// 100,000 rows assigned to it. These will then be loaded into memory, or something else.
//
// When the number of nodes changes, new start and end dates need to be calculated.
//
// When the data is changed, the new end date will need to be accounted for somehow...
//
// We need to make sure that only one node does the rebalancing, otherwise... fucking hell.
// I guess each node as they're aware of which node they are, could just do their bit?
// What if node 3 of 4 goes? Do we need to re-assign node numbers? Probably.
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
