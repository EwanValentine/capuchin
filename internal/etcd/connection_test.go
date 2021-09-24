package etcd

import (
	"context"
	"testing"

	"github.com/EwanValentine/capuchin/conf"
	"github.com/stretchr/testify/require"
)

func TestCanConnect(t *testing.T) {
	conn, err := NewConnection(&conf.Config{
		BrokerAddr: ":2379",
	})
	require.NoError(t, err)
	_, err = conn.KV.Put(context.TODO(), "foo", "bar")
	require.NoError(t, err)
}
