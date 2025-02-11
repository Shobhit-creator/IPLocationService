package middleware

import "fmt"

func NewConcurrentMap() *ConcurrentMap {
    return &ConcurrentMap{
        store: make(map[string]int),
    }
}
func RateLimiter() {
	// 200 requests per second so 1200 requests per minute
	// token bucekt algorithm what it says.
	// a fixed volume of bucket i.e. 100 tokens per bucket
	// refill rate i.e. number of tokens that can be filled in bucket per second i.e. 2 tokens per second

	fmt.Println("RateLimiter")
}