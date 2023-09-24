package ratelimiter

import (
	"math"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/proxy-server-rateLimiter/config"
	"github.com/proxy-server-rateLimiter/database"
)

var timeUnit = 60.00

func rateLimiterInit(endpoint config.Endpoint, key string) {
	rateLimiter := database.RateLimiter{}
	rateLimiter.CurrTime = time.Now()
	rateLimiter.PreCount = float64(endpoint.RequestsPerSec)
	rateLimiter.CurrCount = 1
	database.SetDataInRedis(rateLimiter, key)
}

func RateLimiterHandler(key string, endpoint config.Endpoint) bool {
	rateLimiter, err := database.GetDataFromRedis(key)
	if err == redis.Nil {
		rateLimiterInit(endpoint, key)
	} else {
		timeDiff := time.Since(rateLimiter.CurrTime).Seconds()
		if timeDiff > timeUnit {
			rateLimiter.CurrTime = time.Now()
			rateLimiter.PreCount = rateLimiter.CurrCount
			rateLimiter.CurrCount = 0
		}
		ec := math.Ceil(rateLimiter.PreCount*(timeUnit-timeDiff)/timeUnit + rateLimiter.CurrCount)
		if ec > float64(endpoint.RequestsPerSec) {
			return false
		}

		rateLimiter.CurrCount++
		database.SetDataInRedis(rateLimiter, key)
	}
	return true
}
