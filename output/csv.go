package output

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/vivek2584/faculty-scraper/models"
)

// WriteCSV writes the scraped faculty data to a CSV file at the given path.
func WriteCSV(filename string, faculties []models.Faculty) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(models.CSVHeaders()); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, f := range faculties {
		if err := writer.Write(f.ToRow()); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}
