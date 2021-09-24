package cluster

import (
	"sync"
	"testing"

	"github.com/EwanValentine/capuchin/conf"
	"github.com/EwanValentine/capuchin/internal/etcd"
	"github.com/stretchr/testify/require"
)

var (
	once           sync.Once
	clusterService *Cluster
	start          = 20190101
	end            = 20200101
)

func setUp(t *testing.T) {
	once.Do(func() {
		client, err := etcd.NewConnection(&conf.Config{
			BrokerAddr: "localhost:2379",
		})
		require.NoError(t, err)
		clusterService = NewCluster(client)
	})
}

func TestCanJoinCluster(t *testing.T) {
	setUp(t)
	err := clusterService.Join(0, DataNode, start, end)
	require.NoError(t, err)
}

func TestCanWatchForJoinEvent(t *testing.T) {
	setUp(t)

	clusterService.Watch()

	err := clusterService.Join(0, DataNode, start, end)
	require.NoError(t, err)

	node := <-clusterService.NodeAdded()
	require.Equal(t, start, node.Start)
	require.Equal(t, end, node.End)
}

func TestCanWatchForLeaveEvent(t *testing.T) {
	setUp(t)

	clusterService.Watch()

	err := clusterService.Join(0, DataNode, start, end)
	require.NoError(t, err)
	err = clusterService.Leave(0)
	require.NoError(t, err)

	node := <-clusterService.NodeAdded()
	require.Equal(t, start, node.Start)
	require.Equal(t, end, node.End)
}
