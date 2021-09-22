package cluster

type Status string

var (
	Online  Status = "online"
	Offline Status = "offline"
)

// Node -
type Node struct {
	ID     string
	Status Status
}

// Cluster -
type Cluster struct {
	Nodes []Node
}

// NewCluster -
func NewCluster() *Cluster {
	return &Cluster{}
}

// Heartbeat sends a message to all the other nodes in the cluster
// on a regular interval, any not to respond will be flagged as offline.
// If a node is flagged as offline, it will be removed and will cause
// the shard manager will rebalance. Or something.
func (c *Cluster) Heartbeat() {}
