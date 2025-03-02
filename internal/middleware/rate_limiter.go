package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/shobhit-Creator/IPLocationService/internal/models"
	"github.com/shobhit-Creator/IPLocationService/internal/utils"
	"github.com/spf13/viper"
)

var (
    hourlyBuckets sync.Map
    minuteBuckets sync.Map
    hourlyLimit int
    minuteLimit int
)

func InitRateLimiter() {

    hourlyBuckets = sync.Map{}
    minuteBuckets = sync.Map{}

    go cleanupBuckets(&hourlyBuckets, time.Hour)
    go cleanupBuckets(&minuteBuckets, time.Minute)
} 

func cleanupBuckets(buckets *sync.Map, maxIdleTime time.Duration) {
    for {
        time.Sleep(maxIdleTime)
        now := time.Now()
        buckets.Range(func(key, value interface{}) bool {
            bucket := value.(*models.TokenBucket)
            if now.Sub(bucket.LastRefillTime) > maxIdleTime {
                buckets.Delete(key)
            }
            return true
        })
    }
}

func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
	  ipAddress := utils.GetIPAddress(r)

      hourlyLimit = viper.GetInt("rateLimits.hourlyLimit")
      hourBucket := utils.GetOrCreateBucket(&hourlyBuckets, ipAddress, hourlyLimit, float64(hourlyLimit)/3600)
      if !hourBucket.Allow() {
        http.Error(w, "Service temporarily not available due to hourly rate limiting. Please try again later.", http.StatusTooManyRequests)
        return
      }
      minuteLimit = viper.GetInt("rateLimits.minuteLimit")
      minuteBucket := utils.GetOrCreateBucket(&minuteBuckets, ipAddress, minuteLimit, float64(minuteLimit)/60)
      if !minuteBucket.Allow() {
        http.Error(w, "Service temporarily not available due to minute rate limiting. Please try again later.", http.StatusTooManyRequests)
        return
       }

       next.ServeHTTP(w, r)

    }
}

        // // 200 requests per second so 1200 requests per minute
        // // token bucket algorithm
        // // a fixed volume of bucket i.e. 100 tokens per bucket
        // // refill rate i.e. number of tokens that can be filled in bucket per second i.e. 2 tokens per second

        // // Rate limiting logic here
        // const bucketSize = viper.GetInt("rate_limits.")
		// ipAddress := utils.GetIPAddress(r)
        // currentTime := time.Now().UTC().Format("2006:01:02:15:04")
        // minuteCacheKey := fmt.Sprintf("RateLimit:%s:%s", ipAddress, currentTime)
        // var c interfaces.Cache = cache.GetInstance()
        // val, err := c.Get(minuteCacheKey);
        // var tokenCount int
        // var lastRequestTime time.Time
        // if(err != nil){
        //     tokenCount = 200;
        //     lastRequestTime = time.Now().UTC();
        //     if err := c.Set(minuteCacheKey, []interface{}{tokenCount, lastRequestTime}, time.Minute); err != nil {
        //         http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        //         return
        //     }
        // } else {
        //     data := val.([]interface{})
        //     tokenCount = data[0].(int)
        //     lastRequestTime = data[1].(time.Time)

        //     elapsedTime := time.Now().UTC().Sub(lastRequestTime)

        //     refillTokens := int(elapsedTime.Seconds()) * 2;

        //     tokenCount += refillTokens
        //     if(tokenCount > 200){
        //         tokenCount = 200
        //     }
        //     lastRequestTime = time.Now().UTC()
        // }

        // if tokenCount <= 0 {
        //     w.WriteHeader(http.StatusTooManyRequests);
        //     return
        // }

        // tokenCount -= 1;

        // if err := c.Set(minuteCacheKey, []interface{}{tokenCount, lastRequestTime}, time.Minute); err != nil {
        //     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        //     return
        // }

        // next(w, r)