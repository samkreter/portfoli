package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"io"
	"os"

	"github.com/samkreter/portfoli/reader"
	//"github.com/samkreter/portfoli/allocations"
)

type Percent float64

func main() {

	fileName := "/Users/samkreter/Downloads/Portfolio_Position_Mar-27-2020.csv"

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(lines))

	rows := []*reader.FidelityRow{}
	for _, line := range lines {
		row := reader.ParseRow(line)
		if row != nil {
			rows = append(rows, row)
		}
	}

	for _, row := range rows {
		fmt.Println(row)
	}
}

func test(){

	csvfile, err := os.Open("/Users/samkreter/Downloads/Portfolio_Position_Mar-27-2020.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
		break
	}

	//allocations.GetAllocation("AllWeather")
}
