package orchestration

import (
	"demo/network"
	"demo/service"
	"demo/stack"
)

type EpicSaga struct {
}

func (mode EpicSaga) Orchestrate(network network.Network, orchestrator service.DomainService, receivers []service.DomainService) bool {

	comittedServices := stack.NewStack[service.DomainService]()
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
