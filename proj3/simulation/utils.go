package simulation

import (
	"encoding/csv"
	"os"
	"strconv"
)

type Asset struct {
    Name        string
    InitialPrice float64
    Weight       float64
    MeanReturn   float64
    Volatility   float64
}

func LoadPortfolio(filePath string) ([]Asset, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    var portfolio []Asset
    for i, record := range records {
        if i == 0 { 
            continue
        }
        initialPrice, _ := strconv.ParseFloat(record[1], 64)
        weight, _ := strconv.ParseFloat(record[2], 64)
        meanReturn, _ := strconv.ParseFloat(record[3], 64)
        volatility, _ := strconv.ParseFloat(record[4], 64)

        portfolio = append(portfolio, Asset{
            Name:        record[0],
            InitialPrice: initialPrice,
            Weight:      weight,
            MeanReturn:  meanReturn,
            Volatility:  volatility,
        })
    }
    return portfolio, nil
}
