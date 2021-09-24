package cluster

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	storagepb "github.com/coreos/etcd/storage/storagepb"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

type Status string

var (
	Online  Status = "online"
	Offline Status = "offline"
)

// NodeType -
type NodeType string

var (
	// DataNode is a node that represents a data shard,
	// and exposes access to that shard
	DataNode NodeType = "data_node"
	// ProxyNode queries the data nodes and returns the results
	ProxyNode NodeType = "proxy_node"
)

// Node -
type Node struct {
	ID     int
	Type   NodeType
	Status Status
	Start  int
	End    int
}

// Cluster -
type Cluster struct {
	Nodes   []Node
	client  *clientv3.Client
	added   chan Node
	removed chan Node
	stop    chan struct{}
	errors  chan error
}

// NewCluster -
func NewCluster(client *clientv3.Client) *Cluster {
	return &Cluster{
		client:  client,
		Nodes:   make([]Node, 0),
		added:   make(chan Node),
		removed: make(chan Node),
		stop:    make(chan struct{}),
		errors:  make(chan error),
	}
}

func generateValue(id int, nodeType NodeType, start, end int) string {
	return fmt.Sprintf("%d:%s:%d:%d", id, nodeType, start, end)
}

func parseValue(val string) (*Node, error) {
	parts := strings.Split(val, ":")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing node id, invalid integer")
	}

	// NodeType
	t := parts[1]

	if string(DataNode) == t {
		start, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.New("error parsing start value, invalid integer")
		}

		end, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("error parsing end value, invalid integer")
		}

		return &Node{ID: id, Type: NodeType(t), Start: start, End: end}, nil
	} else if string(ProxyNode) == t {
		return &Node{ID: id, Type: NodeType(t)}, nil
	}

	return nil, errors.New("no valid type found")
}

func generateKey(id int) string {
	return fmt.Sprintf("/nodes/%d", id)
}

// Join -
func (c *Cluster) Join(id int, nodeType NodeType, start, end int) error {
	_, err := c.client.KV.Put(
		context.Background(),
		generateKey(id),
		generateValue(id, nodeType, start, end),
	)
	return err
}

// Leave -
func (c *Cluster) Leave(id int) error {
	_, err := c.client.KV.Delete(context.Background(), generateKey(id))
	return err
}

// Watch -
func (c *Cluster) Watch() {
	go func() {
		eventsChan := c.client.Watch(context.Background(), "/nodes/", clientv3.WithPrefix())
		for {
			select {
			case events := <-eventsChan:
				for _, event := range events.Events {
					val := string(event.Kv.Value)
					if event.Type == storagepb.DELETE {
						log.Println("Node has exited: ", val)

						node, err := parseValue(val)
						if err != nil {
							c.errors <- err
							continue
						}

						c.removed <- *node
					} else if event.Type == storagepb.PUT {
						log.Println("Node has been updated: ", val)

						node, err := parseValue(val)
						if err != nil {
							c.errors <- err
							continue
						}

						c.added <- *node
					}
				}
			case <-c.stop:
				return
			}
		}
	}()
}

// GetNodesByType -
func (c *Cluster) GetNodesByType(nodeType NodeType) ([]*Node, error) {
	nodes := make([]*Node, 0)

	resp, err := c.client.KV.Get(context.Background(), "/nodes/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		node, err := parseValue(string(kv.Value))
		if err != nil {
			return nil, err
		}

		if node.Type == nodeType {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

// NodeRemoved -
func (c *Cluster) NodeRemoved() <-chan Node {
	return c.removed
}

// NodeAdded -
func (c *Cluster) NodeAdded() <-chan Node {
	return c.added
}

// Stop -
func (c *Cluster) Stop() {
	c.stop <- struct{}{}
}
