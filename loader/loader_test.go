package loader

import (
	"testing"
	"time"

	"github.com/EwanValentine/capuchin/cluster"
	"github.com/stretchr/testify/require"
)

type mockDistributionManager struct {
	AddedChan   chan cluster.Node
	RemovedChan chan cluster.Node
}

func (m *mockDistributionManager) NodeAdded() <-chan cluster.Node {
	return m.AddedChan
}

func (m *mockDistributionManager) NodeRemoved() <-chan cluster.Node {
	return m.RemovedChan
}

type mockDataSource struct{}

func (m *mockDataSource) Write(data []byte) error {
	return nil
}

func (m *mockDataSource) Read(start, end int) ([]byte, error) {
	return nil, nil
}

func TestCanRebalanceNodeRemoved(t *testing.T) {
	nodeCount := 4

	dm := &mockDistributionManager{
		AddedChan:   make(chan cluster.Node),
		RemovedChan: make(chan cluster.Node),
	}

	go func() {
		dm.RemovedChan <- cluster.Node{
			ID: 1,
		}
	}()

	ds := &mockDataSource{}

	start := 20190101
	end := 20210920

	periods := GeneratePeriods(nodeCount-1, start, end)

	shardID := 0
	loader := NewLoader(ds, dm, shardID, nodeCount)

	loader.Shard(start, end)
	// Urgh...
	time.Sleep(20 * time.Millisecond)
	loader.Stop()

	require.Equal(t, loader.ShardID, 0)
	require.Equal(t, periods[0].Start, loader.Start)
	require.Equal(t, periods[0].End, loader.End)
}
