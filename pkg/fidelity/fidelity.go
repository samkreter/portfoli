package fidelity

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultFileMatcher = "Portfolio_Position"
)

// GetCurrentPositions parses the passed in Fidelity Portfoli file. If no filepath is used.
// A default filepath of <user's homedir>/Downloads/Protfolio_Positions*. If mutliple files exists, the last modified will be used.
func GetCurrentPositions(filename string) ([]*FidelityRow, error) {
	var err error
	// get default filename
	if filename == "" {
		filename, err = getDefaultFidelityPortfoliPath()
		if err != nil {
			return nil, err
		}
		fmt.Printf("Using default Fidelity portfoli file: %q\n", filename)
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := readCSV(f)
	if err != nil {
		return nil, err
	}

	currentPositions := []*FidelityRow{}
	for _, row := range rows {
		fRow := parseRow(row)
		if fRow != nil {
			currentPositions = append(currentPositions, fRow)
		}
	}

	return currentPositions, nil
}

func readCSV(f *os.File) ([][]string, error) {
	rows := [][]string{}
	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = true

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return rows, nil
			}
			// ignore invalid filed count since the fidelity csv's have junk at the end
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				continue
			}

			return nil, err
		}

		rows = append(rows, row)
	}
}

func getDefaultFidelityPortfoliPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return getLastModifiedFile(usr.HomeDir+"/Downloads", defaultFileMatcher)
}

func getLastModifiedFile(dirPath, substrMatcher string) (string, error) {
	lastModifiedFile := struct {
		path    string
		modTime time.Time
	}{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("ERRROR: %q", err)
			return nil
		}

		// Ensure regular file
		if !info.Mode().IsRegular() {
			return nil
		}

		// contains matching substring
		if !strings.ContainsAny(info.Name(), substrMatcher) {
			return nil
		}

		// Latest modified
		if info.ModTime().Before(lastModifiedFile.modTime) {
			return nil
		}

		lastModifiedFile.path = path
		lastModifiedFile.modTime = info.ModTime()
		return nil
	})
	if err != nil {
		return "", err
	}

	return lastModifiedFile.path, nil
}
