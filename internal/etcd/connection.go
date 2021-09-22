package etcd

import (
	"time"

	"go.etcd.io/etcd/clientv3"
)

// NewConnection -
func NewConnection() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	})
}
