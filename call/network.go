package call

import (
	"math/rand/v2"
	"time"
)

type Network struct {
	AverageTravelLatency time.Duration
	Sigma                time.Duration
}

func (n Network) RandomTravelLatency() time.Duration {
	mean := float64(n.AverageTravelLatency)
	stdDev := float64(n.Sigma)
	latency := rand.NormFloat64()*stdDev + mean
	return time.Duration(latency)
}
