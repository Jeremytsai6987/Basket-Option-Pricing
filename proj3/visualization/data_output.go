package visualization

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func ExportToCSV(data []float64, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write([]string{strconv.FormatFloat(value, 'f', 6, 64)})
		if err != nil {
			log.Fatalf("Failed to write data to CSV: %v", err)
		}
	}
	log.Printf("Data exported to %s\n", filePath)
}
