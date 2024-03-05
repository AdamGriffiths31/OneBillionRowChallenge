package versions

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func RunVersion2(fileName string) string {
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

		value := customParseFloat(valueStr)

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
	return fmt.Sprintf("Running version 2. Time taken: %.2fs", timeInSec)
}

func customParseFloat(valueStr string) float64 {
	// Check if the value is negative
	idx := 0
	negative := false
	if valueStr[idx] == '-' {
		negative = true
		idx++
	}

	// first digit
	value := float64(valueStr[idx] - '0')
	idx++

	// second digit is optional
	if valueStr[idx] != '.' {
		value = value*10 + float64(valueStr[idx]-'0')
		idx++
	}

	// remove decimal point
	idx++

	// decimal part
	value += float64(valueStr[idx]-'0') / 10

	if negative {
		value = -value
	}

	return value
}
