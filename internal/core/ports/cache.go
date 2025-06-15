package ports

import "context"

// CacheService defines the interface for caching operations
type CacheService interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
}
