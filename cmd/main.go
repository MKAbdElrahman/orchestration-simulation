package main

import (
	"demo/call"
	"demo/sweep"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	network := call.Network{
		AverageTravelLatency: 50 * time.Millisecond,
		Sigma:                10 * time.Millisecond,
	}

	alpha := .99
	orchestrator := call.DomainService{
		AverageLatency: 100 * time.Millisecond,
		Sigma:          20 * time.Millisecond,
		Availability:   alpha,
	}

	orchestratedServices := []call.DomainService{
		{
			AverageLatency: 100 * time.Millisecond,
			Sigma:          20 * time.Millisecond,
			Availability:   alpha,
		},
		{
			AverageLatency: 100 * time.Millisecond,
			Sigma:          20 * time.Millisecond,
			Availability:   alpha,
		},
		{
			AverageLatency: 100 * time.Millisecond,
			Sigma:          20 * time.Millisecond,
			Availability:   alpha,
		},
	}
	nonblocking_mode := call.AtomicNonBlockingRequestReply{}
	blocking_mode := call.AtomicBlockingRequestReply{}
	npoints := 100
	sampleSize := 1000
	alpha_0 := .9 
	alpha_f := .999

	experiments := GenerateExperiments(alpha_0, alpha_f, npoints, network, orchestrator, orchestratedServices, nonblocking_mode, sampleSize)

	results := sweep.RunSweep(experiments)

	if err := SaveResultsToCSV(results, "sweep_results_non_blocking.csv"); err != nil {
		fmt.Printf("Error saving results to CSV: %v\n", err)
	} else {
		fmt.Println("Results successfully saved to sweep_results_non_blocking.csv")
	}

	experiments = GenerateExperiments(alpha_0, alpha_f, npoints, network, orchestrator, orchestratedServices, blocking_mode, sampleSize)

	results = sweep.RunSweep(experiments)

	if err := SaveResultsToCSV(results, "sweep_results_blocking.csv"); err != nil {
		fmt.Printf("Error saving results to CSV: %v\n", err)
	} else {
		fmt.Println("Results successfully saved to sweep_results_blocking.csv")
	}

}

// GenerateExperiments generates a slice of experiments with varying availability.
func GenerateExperiments(start, end float64, numPoints int, network call.Network, orchestrator call.DomainService, orchestratedServices []call.DomainService, mode sweep.OrchestratorMode, sampleSize int) []sweep.Experiment {
	step := (end - start) / float64(numPoints-1)
	experiments := make([]sweep.Experiment, numPoints)

	for i := 0; i < numPoints; i++ {
		availability := start + float64(i)*step
		orchestratorCopy := orchestrator
		orchestratorCopy.Availability = availability
		experiments[i] = sweep.Experiment{
			Network:              network,
			Orchestrator:         orchestratorCopy,
			OrchestratorMode:     mode,
			OrchestratedServices: orchestratedServices,
			NumCalls:             sampleSize,
		}
	}

	return experiments
}

func SaveResultsToCSV(results []sweep.Result, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ExperimentIndex", "Availability", "AverageSuccessLatency", "AverageFailureLatency", "SuccessRatio"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for i, result := range results {
		record := []string{
			strconv.Itoa(i),
			fmt.Sprintf("%f", result.Availability),
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
