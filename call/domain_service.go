package call

import (
	"math/rand/v2"
	"time"
)

type DomainService struct {
	AverageLatency time.Duration
	Sigma          time.Duration
	Availability   float64
}

func (s DomainService) RandomWorkLatency() time.Duration {
	mean := float64(s.AverageLatency)
	stdDev := float64(s.Sigma)
	latency := rand.NormFloat64()*stdDev + mean
	return time.Duration(latency)
}

func (s DomainService) IsAvailable() bool {
	return rand.Float64() < s.Availability
}
