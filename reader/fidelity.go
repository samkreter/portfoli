package reader

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	rowLength = 14
	naConst   = "n/a"
	emptyMark = "--"
)

func ParseRow(row []string) *FidelityRow {
	if len(row) < rowLength {
		log.Printf("ERROR: row length is: '%d', expected: '%d'", len(row), rowLength)
		return nil
	}

	var quantity float64
	var err error
	if row[3] != naConst {
		quantity, err = strconv.ParseFloat(row[3], 64)
		if err != nil {
			log.Println(row[3], "ERROR-quantity: ", err)
			return nil
		}
	}

	lastPrice := &Currency{}
	if err := lastPrice.UnmarshalCSV(row[4]); err != nil {
		log.Println("ERROR-lastPrice: ", err)
		return nil
	}

	lastPriceChange := &Currency{}
	err = lastPriceChange.UnmarshalCSV(row[5])
	if err != nil {
		log.Println("ERROR-lastPriceChange: ", err)
		return nil
	}

	current := &Currency{}
	err = current.UnmarshalCSV(row[6])
	if err != nil {
		log.Println("ERROR-current: ", err)
		return nil
	}

	todaysGainLossDollar := &Currency{}
	if err := todaysGainLossDollar.UnmarshalCSV(row[7]); err != nil {
		log.Println("ERROR-todaysGainLossDollar: ", err)
		return nil
	}

	todaysGainLossPercent := Percent(0.0)
	err = todaysGainLossPercent.UnmarshalCSV(row[8])
	if err != nil {
		log.Println(row[8], "   ERROR-todaysGainLossPercent: ", err)
		return nil
	}

	totalGainLossDollar := &Currency{}
	err = totalGainLossDollar.UnmarshalCSV(row[9])
	if err != nil {
		log.Println("ERROR-totalGainLossDollar: ", err)
		return nil
	}

	totalGainLossPercent := Percent(0.0)
	err = totalGainLossPercent.UnmarshalCSV(row[10])
	if err != nil {
		log.Println("ERROR-totalGainLossPercent: ", err)
		return nil
	}

	costBasisPerShare := &Currency{}
	err = costBasisPerShare.UnmarshalCSV(row[11])
	if err != nil {
		log.Println("ERROR-costBasisPerShare: ", err)
		return nil
	}

	costBasisTotal := &Currency{}
	err = costBasisTotal.UnmarshalCSV(row[12])
	if err != nil {
		log.Println("ERROR-costBasisTotal: ", err)
		return nil
	}

	return &FidelityRow{
		AccountName:           row[0],
		Symbol:                row[1],
		Description:           row[2],
		Quantity:              quantity,
		LastPrice:             lastPrice,
		LastPriceChange:       lastPriceChange,
		Current:               current,
		TodaysGainLossDollar:  todaysGainLossDollar,
		TodaysGainLossPercent: todaysGainLossPercent,
		TotalGainLossDollar:   totalGainLossDollar,
		TotalGainLossPercent:  totalGainLossPercent,
		CostBasisPerShare:     costBasisPerShare,
		CostBasisTotal:        costBasisTotal,
		Type:                  row[13],
	}
}

type FidelityRow struct {
	AccountName           string
	Symbol                string
	Description           string
	Quantity              float64
	LastPrice             *Currency
	LastPriceChange       *Currency
	Current               *Currency
	TodaysGainLossDollar  *Currency
	TodaysGainLossPercent Percent
	TotalGainLossDollar   *Currency
	TotalGainLossPercent  Percent
	CostBasisPerShare     *Currency
	CostBasisTotal        *Currency
	Type                  string
}

type Percent float64

func (p Percent) MarshalCSV() (string, error) {
	return fmt.Sprintf("%f", p), nil
}

func (p Percent) String() string {
	return fmt.Sprintf("%f", p)
}

// Convert the CSV string as internal date
func (p Percent) UnmarshalCSV(csv string) (err error) {
	if len(csv) < 2 {
		return nil
	}

	// ignore n/a values
	if csv == naConst {
		return nil
	}

	// ignore the empty markers
	if csv == emptyMark {
		return nil
	}

	// Not a percent
	if string(csv[len(csv)-1]) != "%" {
		return fmt.Errorf("invalid percent")
	}

	valStr := csv[0 : len(csv)-1]

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return err
	}

	p = Percent(val * .001)

	return nil
}

type Currency struct {
	Type  string
	Value float64
}

func (c *Currency) MarshalCSV() (string, error) {
	return fmt.Sprintf("%s%f", c.Type, c.Value), nil
}

func (c *Currency) String() string {
	return fmt.Sprintf("%s%f", c.Type, c.Value) // Redundant, just for example
}

// TODO: Make work with all currency, not just USD
func (c *Currency) UnmarshalCSV(csv string) (err error) {
	if len(csv) < 2 {
		return nil
	}

	// ignore n/a values
	if csv == naConst {
		return nil
	}

	// ignore the empty markers
	if csv == emptyMark {
		return nil
	}

	//c.Type = string(csv[0])
	c.Type = "$"

	valStr := strings.Replace(csv[1:], ",", "", -1)
	valStr = strings.Replace(csv[1:], "$", "", -1)

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return err
	}

	c.Value = val

	return nil
}
