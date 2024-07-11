package orchestration

import (
	"demo/call"
	"demo/stack"
)

type EpicSaga struct {
}

func (mode EpicSaga) Orchestrate(network call.Network, orchestrator call.DomainService, receivers []call.DomainService) bool {

	comittedServices := stack.NewStack[call.DomainService]()
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
