package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/loadimpact/k6/js/modules"
)

func init() {
	modules.Register("k6/x/redis", new(REDIS))
}

// REDIS is the k6 Redis extension.
type REDIS struct{}

// NewClient creates a new Redis client
func (*REDIS) NewClient(addr string, password string, bd int) *redis.Client {
	if addr == "" {
		addr = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       bd,       // use default DB
	})
}

// Set adds a key/value
func (*REDIS) Set(client *redis.Client, key string, value interface{}, expiration time.Duration) {
	// TODO: Make expiration configurable. Or document somewhere the unit.
	err := client.Set(key, value, expiration*time.Second).Err()
	if err != nil {
		ReportError(err, "Failed to set the specified key/value pair")
	}
}

// Get gets a key/value
func (*REDIS) Get(client *redis.Client, key string) string {
	val, err := client.Get(key).Result()
	if err != nil {
		ReportError(err, "Failed to get the specified key")
	}
	return val
}

// Del removes a key/value
func (*REDIS) Del(client *redis.Client, keys ...string) {
	err := client.Del(keys...).Err()
	if err != nil {
		ReportError(err, "Failed to remove the specified keys")
	}
}

// Do runs arbitrary/custom commands
func (*REDIS) Do(client *redis.Client, args ...interface{}) interface{} {
	val, err := client.Do(args...).Result()
	if err != nil {
		if err == redis.Nil {
			ReportError(err, "Key does not exist")
		} else {
			ReportError(err, "Failed to do command")
		}
	}
	return val
}
