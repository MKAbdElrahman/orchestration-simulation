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

type SweepParams struct {
	Network      network.Network        // Network configuration
	Orchestrator service.DomainService  // Orchestrator configuration
	SampleSize   int                    // Sample size for each experiment
	Npoints      int                    // Number of sampling availability points
	Alpha0       float64                // Starting availability for the orchestrated service
	AlphaE       float64                // Ending availability for the orchestrated service
	MinServices  int                    // Minimum number of orchestrated services
	MaxServices  int                    // Maximum number of orchestrated services
	Mode         sweep.OrchestratorMode // Orchestrator mode
	ResultsDir   string                 // Directory to save results
	BaseFileName string                 // Base file name for results
}

func main() {
	orchestratorAvailability := 0.99
	params := SweepParams{
		Network:      configureNetwork(),
		Orchestrator: configureOrchestrator(orchestratorAvailability),
		SampleSize:   1000,
		Npoints:      100,
		Alpha0:       0.900,
		AlphaE:       0.999,
		MinServices:  2,
		MaxServices:  4,
		ResultsDir:   "results",
	}

	modes := []struct {
		mode     sweep.OrchestratorMode
		fileName string
	}{
		{orchestration.FantasyFictionSaga{}, "fantasyfiction"},
		{orchestration.EpicSaga{}, "epic"},
		{orchestration.FairyTaleSaga{}, "fairytale"},
		{orchestration.ParallelSaga{}, "parallel"},
	}

	var wg sync.WaitGroup
	wg.Add(len(modes))

	for _, m := range modes {
		go func(m sweep.OrchestratorMode, fileName string) {
			defer wg.Done()
			params.Mode = m
			params.BaseFileName = fileName
			RunExperimentsAndSaveResults(params)
		}(m.mode, m.fileName)
	}

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

func RunExperimentsAndSaveResults(params SweepParams) {
	var wg sync.WaitGroup

	for numServices := params.MinServices; numServices <= params.MaxServices; numServices++ {
		wg.Add(1)
		go func(numServices int) {
			defer wg.Done()
			outputFileName := fmt.Sprintf("%s_%d_services.csv", params.BaseFileName, numServices)
			experiments := generateExperiments(params, numServices)
			results := sweep.RunExperiments(experiments)
			if err := saveResultsToCSV(results, params.ResultsDir, outputFileName); err != nil {
				fmt.Printf("Error saving results to CSV: %v\n", err)
			} else {
				fmt.Printf("Results successfully saved to %s\n", outputFileName)
			}
		}(numServices)
	}

	wg.Wait()
}

func generateExperiments(params SweepParams, numServices int) []sweep.Experiment {
	step := (params.AlphaE - params.Alpha0) / float64(params.Npoints-1)
	experiments := make([]sweep.Experiment, 0, params.Npoints)

	for i := 0; i < params.Npoints; i++ {
		availability := params.Alpha0 + float64(i)*step
		orchestratorCopy := params.Orchestrator
		orchestratorCopy.Availability = availability

		orchestratedServicesCopy := configureOrchestratedServices(availability, numServices)
		experiments = append(experiments, sweep.Experiment{
			Network:              params.Network,
			Orchestrator:         orchestratorCopy,
			OrchestratorMode:     params.Mode,
			OrchestratedServices: orchestratedServicesCopy,
			NumCalls:             params.SampleSize,
		})
	}

	return experiments
}
