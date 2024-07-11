package network

import (
	"demo/service"
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

func (network Network) Call(sender service.DomainService, receiver service.DomainService) bool {
	if !sender.IsAvailable() {
		return false
	}
	// travel to receiver
	time.Sleep(network.RandomTravelLatency())

	if !receiver.IsAvailable() {
		return false
	}
	time.Sleep(receiver.RandomWorkLatency())

	// travel back to sender
	time.Sleep(network.RandomTravelLatency())

	// sender has to be available to receive the request
	return sender.IsAvailable()
}
