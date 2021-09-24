package data

import (
	"github.com/EwanValentine/capuchin/cluster"
	"github.com/EwanValentine/capuchin/conf"
	"github.com/EwanValentine/capuchin/internal/etcd"
	"github.com/EwanValentine/capuchin/loader"
	"github.com/EwanValentine/capuchin/source"
)

// Service -
type Service struct{}

// NewService -
func NewService() (*Service, error) {
	c := conf.Load()
	client, err := etcd.NewConnection(c)
	if err != nil {
		return nil, err
	}

	clusterManager := cluster.NewCluster(client)
	existingNodes, err := clusterManager.GetNodesByType(cluster.DataNode)
	if err != nil {
		return nil, err
	}

	// We could maybe just use a counter in etcd for this purpose,
	// would likely be slightly quicker
	nodeID := len(existingNodes)

	// We need to work out the start and end date of the dataset, we can just
	// pass it in for now?

	start := 20190101
	end := 20210101

	// Generates a list of periods for the given data range by date
	// gets the start and end date for this node's assigned period
	periods := loader.GeneratePeriods(nodeID, start, end)
	newStart := periods[nodeID].Start
	newEnd := periods[nodeID].End

	// Join the cluster
	clusterManager.Join(nodeID, cluster.DataNode, newStart, newEnd)

	// Listen for changes
	clusterManager.Watch()

	dataSource := source.NewFileSource("./test-data/")
	loader.NewLoader(dataSource, clusterManager, start, end)

	return &Service{}, nil
}
