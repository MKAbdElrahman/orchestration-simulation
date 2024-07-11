package main

import (
	"demo/cmd/sweep"
	"demo/network"
	"demo/orchestration"
	"demo/service"
	"fmt"
	"sync"
	"time"
)

func main() {
	network := configureNetwork()

	orchestrator_availability := .99
	orchestrator := configureOrchestrator(orchestrator_availability)
	// each measurement is averaged over a sample size due to the probabilistic nature distibuted systems
	sampleSize := 1000
	alpha_0 := .9                    // starting availability for the orchestrated service
	alpha_e := .999                  // ending availability for the orchestrated service
	npoints := 100                   //  number or sampling availability points
	minServices, maxServices := 2, 4 // how many orchestrated services
	resultsDir := "results"

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alpha_0, alpha_e, minServices, maxServices, orchestration.FantasyFictionSaga{}, resultsDir, "fantasyfiction")
	}()

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alpha_0, alpha_e, minServices, maxServices, orchestration.EpicSaga{}, resultsDir, "epic")
	}()

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alpha_0, alpha_e, minServices, maxServices, orchestration.FairyTaleSaga{}, resultsDir, "fairytale")
	}()

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alpha_0, alpha_e, minServices, maxServices, orchestration.ParallelSaga{}, resultsDir, "parallel")
	}()

	wg.Wait()
	fmt.Println("All sweeps completed")
}

func configureNetwork() network.Network {
	return network.Network{
		AverageTravelLatency: 50 * time.Millisecond,
		Sigma:                10 * time.Millisecond,
	}
}

func configureOrchestrator(availability float64) service.DomainService {
	return service.DomainService{
		AverageLatency: 100 * time.Millisecond,
		Sigma:          20 * time.Millisecond,
		Availability:   availability,
	}
}

func configureOrchestratedServices(availability float64, numServices int) []service.DomainService {
	services := make([]service.DomainService, numServices)
	for i := range services {
		services[i] = service.DomainService{
			AverageLatency: 100 * time.Millisecond,
			Sigma:          20 * time.Millisecond,
			Availability:   availability,
		}
	}
	return services
}

func runSweep(network network.Network, orchestrator service.DomainService, sampleSize, npoints int, alpha_0, alpha_e float64, minServices, maxServices int, mode sweep.OrchestratorMode, resultsDir string, baseFileName string) {

	var wg sync.WaitGroup

	for numServices := minServices; numServices <= maxServices; numServices++ {
		wg.Add(1)
		go func(numServices int) {
			defer wg.Done()
			outputFileName := fmt.Sprintf("%s_%d_services.csv", baseFileName, numServices)
			experiments := generateExperiments(alpha_0, alpha_e, npoints, network, orchestrator, mode, sampleSize, numServices)
			results := sweep.RunSweep(experiments)
			if err := saveResultsToCSV(results, resultsDir, outputFileName); err != nil {
				fmt.Printf("Error saving results to CSV: %v\n", err)
			} else {
				fmt.Printf("Results successfully saved to %s\n", outputFileName)
			}
		}(numServices)
	}

	wg.Wait()
}

func generateExperiments(start, end float64, numPoints int, network network.Network, orchestrator service.DomainService, mode sweep.OrchestratorMode, sampleSize, numServices int) []sweep.Experiment {
	step := (end - start) / float64(numPoints-1)
	experiments := make([]sweep.Experiment, 0, numPoints)

	for i := 0; i < numPoints; i++ {
		availability := start + float64(i)*step
		orchestratorCopy := orchestrator
		orchestratorCopy.Availability = availability

		orchestratedServicesCopy := configureOrchestratedServices(availability, numServices)
		experiments = append(experiments, sweep.Experiment{
			Network:              network,
			Orchestrator:         orchestratorCopy,
			OrchestratorMode:     mode,
			OrchestratedServices: orchestratedServicesCopy,
			NumCalls:             sampleSize,
		})
	}

	return experiments
}
