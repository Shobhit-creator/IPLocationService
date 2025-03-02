package interfaces

import "time"

type Cache interface {
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}
