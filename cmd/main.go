package main

import (
	"demo/cmd/sweep"
	"demo/network"
	"demo/orchestration"
	"demo/service"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func main() {
	network := configureNetwork()
	orchestrator := configureOrchestrator(.99)
	sampleSize := 1000
	npoints := 100
	alphaStart, alphaEnd := .9, .999
	minServices, maxServices := 2, 4

	var wg sync.WaitGroup

	wg.Add(4) // We are running four sweeps in parallel

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alphaStart, alphaEnd, minServices, maxServices, orchestration.FantasyFictionSaga{}, "fantasyfiction")
	}()

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alphaStart, alphaEnd, minServices, maxServices, orchestration.EpicSaga{}, "epic")
	}()

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alphaStart, alphaEnd, minServices, maxServices, orchestration.FairyTaleSaga{}, "fairytale")
	}()

	go func() {
		defer wg.Done()
		runSweep(network, orchestrator, sampleSize, npoints, alphaStart, alphaEnd, minServices, maxServices, orchestration.ParallelSaga{}, "parallel")
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

func runSweep(network network.Network, orchestrator service.DomainService, sampleSize, npoints int, alphaStart, alphaEnd float64, minServices, maxServices int, mode sweep.OrchestratorMode, baseFileName string) {
	var wg sync.WaitGroup

	for numServices := minServices; numServices <= maxServices; numServices++ {
		wg.Add(1)
		go func(numServices int) {
			defer wg.Done()
			outputFileName := fmt.Sprintf("%s_%d_services.csv", baseFileName, numServices)
			experiments := generateExperiments(alphaStart, alphaEnd, npoints, network, orchestrator, mode, sampleSize, numServices)
			results := sweep.RunSweep(experiments)
			if err := saveResultsToCSV(results, outputFileName); err != nil {
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

func saveResultsToCSV(results []sweep.Result, fileName string) error {
	resultsDir := "results"

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(resultsDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(resultsDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ExperimentIndex", "Availability", "NumDomainServices", "AverageSuccessLatency", "AverageFailureLatency", "SuccessRatio"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for i, result := range results {
		record := []string{
			strconv.Itoa(i),
			fmt.Sprintf("%f", result.Availability),
			strconv.Itoa(result.NumDomainServices),
			result.AverageSuccessLatency.String(),
			result.AverageFailureLatency.String(),
			fmt.Sprintf("%.2f%%", result.SuccessRatio*100),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
