package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vivek2584/faculty-scraper/output"
	"github.com/vivek2584/faculty-scraper/scraper"
)

func main() {
	outputFile := flag.String("o", "faculty_data.csv", "output CSV file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <department-or-faculty-URL>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  # Scrape all faculty from a department page:")
		fmt.Fprintln(os.Stderr, "  go run main.go https://www.srmist.edu.in/department/college-of-occupational-therapy/")
		fmt.Fprintln(os.Stderr, "\n  # Scrape a single faculty profile to a custom file:")
		fmt.Fprintln(os.Stderr, "  go run main.go -o results.csv https://www.srmist.edu.in/faculty/dr-ganapathy-sankar-u/")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	startURL := flag.Arg(0)

	faculties, err := scraper.Scrape(startURL)
	if err != nil {
		log.Printf("Warning: %v", err)
	}

	if len(faculties) == 0 {
		fmt.Println("\nNo faculty data found. Try a different URL.")
		return
	}

	if err := output.WriteCSV(*outputFile, faculties); err != nil {
		log.Fatalf("Failed to write CSV: %v", err)
	}

	fmt.Printf("\nDone! Scraped %d faculty members → %s\n", len(faculties), *outputFile)
}
