package main

import (
	"adamgriffiths/1brc/versions"
	"flag"
	"fmt"
)

func main() {
	version := flag.String("version", "1", "Which version of the app to run")
	fileName := flag.String("file", "testdata", "The file to read data from")
	flag.Parse()

	fmt.Printf("Running version %s with file %s\n", *version, *fileName)
	switch *version {
	case "1":
		fmt.Println(versions.RunVersion1(*fileName))
	case "2":
		fmt.Println(versions.RunVersion2(*fileName))
	default:
		fmt.Println("Invalid version")
	}
}
