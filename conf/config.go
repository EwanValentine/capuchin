package conf

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

// Load -
func load(c interface{}) error {
	ctx := context.Background()
	return envconfig.Process(ctx, c)
}

// Load -
func Load() *Config {
	var c Config
	if err := load(&c); err != nil {
		log.Panic(err)
	}
	return &c
}

// Config -
type Config struct {
	HostAddr string `env:"HOST_ADDR,default=:8080"`
	GRPCAddr string `env:"GRPC_ADDR,default=:9090"`
}
