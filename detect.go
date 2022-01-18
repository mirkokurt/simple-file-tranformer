package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/therecipe/qt/widgets"
)

func isLoytech(inputfile string) (bool, error) {
	// Open the input file
	fr, err := os.Open(inputfile)
	if err != nil {
		return false, err
	}
	defer fr.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(fr)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'

	rec, err := csvReader.Read()

	// If the input file is not a LOYTEC, return true, false otherwise
	if err == io.EOF || !strings.Contains(rec[0], "LOYTEC") {
		return false, nil
	}

	return true, nil
}

func transformLoytecToModulo5(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) (string, error) {
	// Open the input file
	fr, err := os.Open(inputfile)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(fr)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'

	// If the input file is a LOYTEC execute a pre-conversion in Modulo 5 format
	// Create a temp output file
	fo, err := os.Create(outputfile + ".temp")
	if err != nil {
		return "", errors.New("attenzione il file di output è aperto e non può essere scritto")
	}
	// Close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Printf("Error closing output file: %s", err)
		}
	}()

	// Create a writer
	writer := csv.NewWriter(fo)
	defer writer.Flush()

	writer.Comma = ';'

	// Init the line counter
	line := 0
	channelNumber := 0
	for {
		// Increment the line
		line++
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		checkError("Cannot read the file", err)

		// If one of the first two lines
		if line < 3 {
			empty_rec := []string{"", "", ""}
			err = writer.Write(empty_rec)
			checkError("Cannot write to file", err)
		} else if line >= 3 && line < 7 {
			// Skip the lines
		} else if line == 7 {
			// Write the header of Modulo 5
			err = writer.Write(table_headers["modulo5_old"])
			checkError("Cannot write to file", err)
		} else {
			// Write the scrambled line
			err = rewriteFromLoytecToModulo5(writer, rec, channelNumber)
			checkError("Cannot write to file", err)
			channelNumber++
		}
		if err != nil {
			return "", err
		}
	}

	return outputfile + ".temp", nil

}

func idOldModulo5(check string) bool {
	if check == "triggerd" || check == "Trigger: mode" {
		return true
	}
	return false
}
