package versions

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

func RunVersion3(fileName string) string {
	startTime := time.Now()

	file, err := os.Open(fileName + ".txt")
	if err != nil {
		return fmt.Sprintf("Error opening file: %s", err)
	}
	defer file.Close()
	startPointer := 0
	cityMap := make(map[string]Data, 1000000)
	buffer := make([]byte, 1024*1024)
	for {
		bytesRead, err := file.Read(buffer[startPointer:])
		if err != nil && err != io.EOF {
			break
		}

		if startPointer+bytesRead == 0 {
			fmt.Println("No bytes read")
			break
		}

		fileData := buffer[:startPointer+bytesRead]

		newLineIndex := bytes.LastIndexByte(fileData, '\n')
		if newLineIndex < 0 {
			fmt.Println("No new line found for: ", string(fileData))
			break
		}

		remainingFileData := fileData[newLineIndex+1:]
		fileData = fileData[:newLineIndex+1]

		for {
			cityData, remainingData, ok := bytes.Cut(fileData, []byte(";"))
			if !ok {
				break
			}
			if len(remainingData) < 4 {
				break
			}

			// Check if the value is negative
			idx := 0
			isNegative := false
			if remainingData[idx] == '-' {
				isNegative = true
				idx++
			}

			// first digit
			value := float64(remainingData[idx] - '0')
			idx++

			// second digit is optional
			if remainingData[idx] != '.' {
				value = value*10 + float64(remainingData[idx]-'0')
				idx++
			}

			// remove decimal point
			idx++

			// decimal part
			value += float64(remainingData[idx]-'0') / 10
			if isNegative {
				value = -value
			}
			idx += 2

			fileData = remainingData[idx:]
			data, ok := cityMap[string(cityData)]
			if !ok {
				data = Data{Count: 0, Min: value, Mean: 0, Max: value, Sum: 0, City: string(cityData)}
			}

			data.Count++
			if value < data.Min {
				data.Min = value
			}
			if value > data.Max {
				data.Max = value
			}
			data.Sum += value
			cityMap[string(cityData)] = data
		}
		startPointer = copy(buffer, remainingFileData)
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
	return fmt.Sprintf("Running version 3. Time taken: %.2fs", timeInSec)
}
