package call

import (
	"demo/stack"
	"sync"
)

type AtomicBlockingRequestReply struct {
}
type AtomicNonBlockingRequestReply struct {
}

func (mode AtomicBlockingRequestReply) Orchestrate(network Network, orchestrator DomainService, receivers []DomainService) bool {

	comittedServices := stack.NewStack[DomainService]()
	for _, svc := range receivers {
		if network.Call(orchestrator, svc) {
			comittedServices.Push(svc)
			if comittedServices.Size() == len(receivers) {
				return true
			}
			continue
		} else {
			break
		}
	}
	// compensate
	for comittedServices.Size() != 0 {
		svc, _ := comittedServices.Pop()
		if network.Call(orchestrator, svc) {
			continue
		} else { // retry
			comittedServices.Push(svc)
		}
	}
	return false
}

func (mode AtomicNonBlockingRequestReply) Orchestrate(network Network, orchestrator DomainService, receivers []DomainService) bool {
	committedServicesCh := make(chan DomainService, len(receivers))
	var wg sync.WaitGroup
	for _, svc := range receivers {
		wg.Add(1)
		go func(svc DomainService) {
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
		go func(svc DomainService) {
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
