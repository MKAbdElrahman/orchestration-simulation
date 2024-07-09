package call

import "time"

func (network Network) Call(sender DomainService, receiver DomainService) bool {
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
