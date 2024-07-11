package orchestration

import (
	"demo/call"
	"sync"
)

type FantasyFictionSaga struct {
}

func (mode FantasyFictionSaga) Orchestrate(network call.Network, orchestrator call.DomainService, receivers []call.DomainService) bool {
	committedServicesCh := make(chan call.DomainService, len(receivers))
	var wg sync.WaitGroup
	for _, svc := range receivers {
		wg.Add(1)
		go func(svc call.DomainService) {
			defer wg.Done()
			if network.Call(orchestrator, svc) {
				committedServicesCh <- svc
			}
		}(svc)
	}
	wg.Wait()

	if len(committedServicesCh) == len(receivers) {
		close(committedServicesCh)
		return true
	}

	for len(committedServicesCh) != 0 {
		wg.Add(1)
		svc := <-committedServicesCh
		go func(svc call.DomainService) {
			defer wg.Done()
			if !network.Call(orchestrator, svc) {
				committedServicesCh <- svc
			}
		}(svc)
	}
	wg.Wait()
	close(committedServicesCh)
	return false
}
