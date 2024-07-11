package orchestration

import (
	"demo/network"
	"demo/service"
	"sync"
)

type ParallelSaga struct {
}

func (mode ParallelSaga) Orchestrate(network network.Network, orchestrator service.DomainService, receivers []service.DomainService) bool {
	committedServicesCh := make(chan service.DomainService, len(receivers))
	var wg sync.WaitGroup
	for _, svc := range receivers {
		wg.Add(1)
		go func(svc service.DomainService) {
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
	// compensate in the background
	go func() {
		for len(committedServicesCh) != 0 {
			wg.Add(1)
			svc := <-committedServicesCh
			go func(svc service.DomainService) {
				defer wg.Done()
				if !network.Call(orchestrator, svc) {
					committedServicesCh <- svc
				}
			}(svc)
		}
		wg.Wait()
		close(committedServicesCh)
	}()

	return false

}
