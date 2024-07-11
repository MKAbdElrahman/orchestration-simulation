package main

import (
	"demo/cmd/sweep"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func saveResultsToCSV(results []sweep.Result, resultsDir string, fileName string) error {
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
