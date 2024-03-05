package versions

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RunVersion1(fileName string) string {
	startTime := time.Now()

	file, err := os.Open(fileName + ".txt")
	if err != nil {
		return fmt.Sprintf("Error opening file: %s", err)
	}
	defer file.Close()

	cityMap := make(map[string]Data, 1000000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		city, valueStr, ok := strings.Cut(line, ";")
		if !ok {
			return fmt.Sprintf("Error parsing line: %s", line)
		}

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return fmt.Sprintf("Error parsing value: %s", err)
		}

		data, ok := cityMap[city]
		if !ok {
			data = Data{Count: 0, Min: value, Mean: 0, Max: value, Sum: 0, City: city}
		}

		data.Count++
		if value < data.Min {
			data.Min = value
		}
		if value > data.Max {
			data.Max = value
		}
		data.Sum += value
		cityMap[city] = data
	}

	sortedData := make([]Data, 0, len(cityMap))
	for _, data := range cityMap {
		data.Mean = data.Sum / float64(data.Count)
		sortedData = append(sortedData, data)
	}

	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].City < sortedData[j].City
	})

	outputFile, err := os.Create("go_output.txt")
	if err != nil {
		return fmt.Sprintf("Error creating file: %s", err)
	}
	defer outputFile.Close()

	fmt.Fprint(outputFile, "{")
	for i, data := range sortedData {
		fmt.Fprintf(outputFile, "%s=%.2f/%.2f/%.2f", data.City, data.Min, data.Mean, data.Max)
		if i < len(sortedData)-1 {
			fmt.Fprint(outputFile, ", ")
		}
	}
	fmt.Fprint(outputFile, "}")

	timeInSec := time.Since(startTime).Seconds()
	return fmt.Sprintf("Running version 1. Time taken: %.2fs", timeInSec)
}
