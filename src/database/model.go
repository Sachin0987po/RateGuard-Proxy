package database

import (
	"time"
)

type RateLimiter struct {
	CurrTime  time.Time `json:"currTime"`
	CurrCount float64   `json:"currCount"`
	PreCount  float64   `json:"preCount"`
}