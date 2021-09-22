package cluster

import (
	"context"
	"fmt"
	"log"
	"time"

	storagepb "github.com/coreos/etcd/storage/storagepb"
	"go.etcd.io/etcd/clientv3"
)

type Status string

var (
	Online  Status = "online"
	Offline Status = "offline"
)

// Node -
type Node struct {
	ID     int
	Status Status
	Start  time.Time
	End    time.Time
}

// Cluster -
type Cluster struct {
	Nodes  []Node
	client *clientv3.Client
}

// NewCluster -
func NewCluster(client *clientv3.Client) *Cluster {
	return &Cluster{client: client, Nodes: make([]Node, 0)}
}

// Heartbeat sends a message to all the other nodes in the cluster
// on a regular interval, any not to respond will be flagged as offline.
// If a node is flagged as offline, it will be removed and will cause
// the shard manager will rebalance. Or something.
func (c *Cluster) Heartbeat() {}

// NodeType -
type NodeType string

var (
	DataNode  NodeType = "data_node"
	ProxyNode NodeType = "proxy_node"
)

func generateValue(id string, nodeType NodeType) string {
	return fmt.Sprintf("%s:%s", id, nodeType)
}

// Join -
func (c *Cluster) Join(id string, nodeType NodeType) error {
	_, err := c.client.KV.Put(context.Background(), "/nodes/"+id, generateValue(id, nodeType))
	return err
}

// Watch -
func (c *Cluster) Watch() {
	go func() {
		eventsChan := c.client.Watch(context.Background(), "/nodes/", clientv3.WithPrefix())
		for events := range eventsChan {
			for _, event := range events.Events {
				val := event.Kv.Value
				if event.Type == storagepb.DELETE {
					log.Println("Node has exited: ", val)
				} else if event.Type == storagepb.PUT {
					log.Println("Node has been updated: ", val)
				}
			}
		}
	}()
}
