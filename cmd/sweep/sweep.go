package sweep

import (
	"demo/network"
	"demo/service"
	"sync"
	"time"
)

type OrchestratorMode interface {
	Orchestrate(ds network.Network, orchestrator service.DomainService, receivers []service.DomainService) bool
}

type Experiment struct {
	Network              network.Network
	OrchestratorMode     OrchestratorMode
	Orchestrator         service.DomainService
	OrchestratedServices []service.DomainService
	NumCalls             int
}

// Result represents the result of an experiment.
type Result struct {
	Availability          float64
	AverageSuccessLatency time.Duration
	AverageFailureLatency time.Duration
	SuccessRatio          float64
	NumDomainServices     int
}

// RunExperiment runs a single experiment and returns the result.
func RunExperiment(exp Experiment) Result {
	successLatencyChan := make(chan time.Duration, exp.NumCalls)
	failureLatencyChan := make(chan time.Duration, exp.NumCalls)
	var wg sync.WaitGroup
	var failCount int
	var successCount int
	var mu sync.Mutex // Mutex to protect failCount and successCount

	// Perform services concurrently
	for i := 0; i < exp.NumCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startTime := time.Now()
			ok := exp.OrchestratorMode.Orchestrate(exp.Network, exp.Orchestrator, exp.OrchestratedServices)
			latency := time.Since(startTime)

			mu.Lock()
			if ok {
				successCount++
				successLatencyChan <- latency
			} else {
				failCount++
				failureLatencyChan <- latency
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
	close(successLatencyChan)
	close(failureLatencyChan)

	totalSuccessLatency := time.Duration(0)
	for latency := range successLatencyChan {
		totalSuccessLatency += latency
	}

	totalFailureLatency := time.Duration(0)
	for latency := range failureLatencyChan {
		totalFailureLatency += latency
	}

	averageSuccessLatency := time.Duration(0)
	if successCount > 0 {
		averageSuccessLatency = time.Duration(totalSuccessLatency.Nanoseconds() / int64(successCount))
	}

	averageFailureLatency := time.Duration(0)
	if failCount > 0 {
		averageFailureLatency = time.Duration(totalFailureLatency.Nanoseconds() / int64(failCount))
	}

	failureRatio := float64(failCount) / float64(exp.NumCalls)
	successRatio := 1.0 - failureRatio

	return Result{
		Availability:          exp.Orchestrator.Availability,
		AverageSuccessLatency: averageSuccessLatency,
		AverageFailureLatency: averageFailureLatency,
		SuccessRatio:          successRatio,
		NumDomainServices:     len(exp.OrchestratedServices),
	}
}

// RunSweep runs a sweep of experiments with different configurations.
func RunExperiments(experiments []Experiment) []Result {
	results := make([]Result, len(experiments))

	for i, exp := range experiments {
		results[i] = RunExperiment(exp)
	}

	return results
}
