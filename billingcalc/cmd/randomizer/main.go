package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	fpFlagValue = flag.String("f", "", "CSV filepath to randomize")
	debug       = flag.Bool("d", false, "Enable debug info")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: ./randomizer [flags] <output file>") //nolint:forbidigo
		os.Exit(1)
	}

	if !*debug {
		log.SetOutput(io.Discard)
	}

	buildNames()

	// Get csv filepath
	fp := *fpFlagValue
	if fp == "" {
		fmt.Print("Enter the CSV filepath to randomize: ") //nolint:forbidigo
		_, err := fmt.Scanln(&fp)
		if err != nil {
			log.Panic(err)
		}
	}

	// Parse csv
	csvFile, err := os.Open(fp)
	if err != nil {
		fmt.Printf("Invalid file: %s\n", fp) //nolint:forbidigo
		os.Exit(1)
	}
	log.Printf("Found %s\n", fp)
	csvTable, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Panic(err)
	}

	// Get columns to randomize
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose which columns to randomize (A-Z): ") //nolint:forbidigo
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	line = strings.ToUpper(line)
	line = strings.Trim(line, "\n")

	colSlice := strings.Split(line, " ")
	log.Printf("Randomizing columns %v\n", colSlice)
	colPositions := make([]int, len(colSlice))
	for i, col := range colSlice {
		colPositions[i] = getColumnPosition(col)
	}
	log.Printf("Column positions %v\n", colPositions)

	// Randomize each column
	knownNames := make(map[string]string)
	for _, row := range csvTable {
		for _, col := range colPositions {
			cell := row[col]
			cell = strings.ToUpper(cell)
			if _, ok := knownNames[cell]; !ok {
				name, err := getRandomName()
				if err != nil {
					log.Panic(err)
				}
				knownNames[cell] = name
				log.Printf("Replacing %s with %s", cell, knownNames[cell])
			}
			row[col] = knownNames[cell]
		}
	}

	// Write output
	outFilepath := flag.Arg(0)
	log.Printf("Using %s as the output file.", outFilepath)
	outFile, err := os.Create(outFilepath)
	if err != nil {
		log.Panic(err)
	}

	if err := csv.NewWriter(outFile).WriteAll(csvTable); err != nil {
		log.Panic(err)
	}
	log.Printf("Remaining names: %d", len(names))

	fmt.Println("Success!") //nolint:forbidigo
}

func getColumnPosition(col string) int {
	return int(col[0] - 65) // ASCII A = 65
}
