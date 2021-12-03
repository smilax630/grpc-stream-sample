package gateway

import (
	"time"

	"github.com/go-redis/redis"
)

const (
	expiration = time.Minute * 10
)

// NewRedisClient Initiates a new client
func NewRedisClient(opts RedisOption) (*redis.Client, error) {

	client := redis.NewClient(opts.ClientOptions())

	// ping
	_, err := client.Ping().Result()
	if err != nil {
		return nil, toGatewayError(err)
	}

	return client, nil
}

// RedisOption is redis option settings
type RedisOption struct {
	// Host address with port number
	// For normal client will only used the first value
	Hosts string

	// Database to be selected after connecting to the server.
	Database int
	// Automatically adds a prefix to all keys
	KeyPrefix string

	// Following options are copied from Options struct.
	Password string

	// Default is 5 seconds.
	DialTimeout time.Duration

	// Default is 10 connections.
	PoolSize int
}

// ClientOptions translates current configuration into a *redis.Options struct
func (o *RedisOption) ClientOptions() *redis.Options {
	opts := &redis.Options{
		Addr:        o.Hosts,
		Password:    o.Password,
		DB:          o.Database,
		DialTimeout: o.DialTimeout,
		PoolSize:    o.PoolSize,
	}
	return opts
}
