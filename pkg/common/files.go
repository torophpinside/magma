package common

import (
	"bufio"
	"encoding/csv"
	"os"
)

func FileReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func SaveToCSV(data [][]string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvwriter := csv.NewWriter(file)
	defer csvwriter.Flush()

	for _, row := range data {
		csvwriter.Write(row)
	}

	return nil
}
