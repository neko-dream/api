package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {
	var (
		csvPath    = flag.String("csv", "", "Path to CSV file containing TalkSession UUIDs")
		outputPath = flag.String("output", "", "Output SQL file path (optional, defaults to stdout)")
	)
	flag.Parse()

	if *csvPath == "" {
		log.Fatal("CSV file path is required. Use -csv flag")
	}

	// Open CSV file
	file, err := os.Open(*csvPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Read CSV
	reader := csv.NewReader(file)

	// Read header if exists
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read CSV header: %v", err)
	}

	// Find UUID column index
	uuidColIndex := -1
	for i, col := range header {
		colLower := strings.ToLower(strings.TrimSpace(col))
		if colLower == "uuid" || colLower == "talk_session_id" || colLower == "talksessionid" {
			uuidColIndex = i
			break
		}
	}

	if uuidColIndex == -1 {
		// If no header matches, assume first column is UUID
		uuidColIndex = 0
		// Reset reader to beginning
		_, err := file.Seek(0, 0)
		if err != nil {
			log.Fatalf("Failed to seek file: %v", err)
		}
		reader = csv.NewReader(file)
	}

	// Prepare output
	var output io.Writer = os.Stdout
	if *outputPath != "" {
		outFile, err := os.Create(*outputPath)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer outFile.Close()
		output = outFile
	}

	// Write SQL header
	fmt.Fprintln(output, "-- Generated INSERT statements for talk_session_report_histories")
	fmt.Fprintf(output, "-- Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(output, "-- CSV file: %s\n\n", *csvPath)

	// Process CSV rows
	insertCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			continue
		}

		if uuidColIndex >= len(record) {
			log.Printf("Skipping row: not enough columns")
			continue
		}

		uuids := strings.TrimSpace(record[uuidColIndex])
		if uuids == "" || uuids == "uuid" || uuids == "talk_session_id" {
			// Skip empty or header rows
			continue
		}

		// Validate UUID format (basic check)
		if len(uuids) != 36 || strings.Count(uuids, "-") != 4 {
			log.Printf("Warning: Invalid UUID format: %s", uuids)
			continue
		}

		// Generate UUIDv7 for the history ID
		historyID, err := uuid.NewV7()
		if err != nil {
			log.Printf("Failed to generate UUIDv7: %v", err)
			continue
		}

		// Generate INSERT statement with SELECT
		fmt.Fprintf(output, `INSERT INTO talk_session_report_histories (talk_session_report_history_id, talk_session_id, report, created_at)
SELECT '%s'::uuid, talk_session_id, report, created_at
FROM talk_session_reports
WHERE talk_session_id = '%s'::uuid;
`,
			historyID.String(),
			uuids,
		)
		insertCount++
	}

	fmt.Fprintf(output, "\n-- Total INSERT statements: %d\n", insertCount)

	if *outputPath != "" {
		fmt.Printf("Generated %d INSERT statements in %s\n", insertCount, *outputPath)
	} else {
		fmt.Fprintf(os.Stderr, "\nGenerated %d INSERT statements\n", insertCount)
	}
}