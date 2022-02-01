package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/therecipe/qt/widgets"
)

func transformFromModulo5ToSelected(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) error {

	// Open the input file
	fr, err := os.Open(inputfile)
	if err != nil {
		return err
	}
	defer fr.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(fr)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'

	// Open the output file
	fo, err := os.Create(outputfile)
	if err != nil {
		return errors.New("attenzione il file di output è aperto e non può essere scritto")
	}
	// Close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Printf("Error closing output file: ", err)
		}
	}()

	// Create a writer
	writer := csv.NewWriter(fo)
	defer writer.Flush()

	writer.Comma = ';'

	// Init the line counter
	line := 0
	isOld := false

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
			// TODO: For "modulo 6" skip the first lines Understand why the "modulo6" keeps stopping import.  rec[0] = "EY6AS80 Modbus Master" does not fully work
			if selectedOutput != "modulo6" {
				// Write the first two line as is (exept the first for "modulo6")
				err = writer.Write(rec)
				checkError("Cannot write to file", err)
			}

		} else if line == 3 {
			//Detect if it is an old version or a new one
			isOld = idOldModulo5(rec[10])

			// Write the selected device table header
			err = writer.Write(table_headers[selectedOutput])
			checkError("Cannot write to file", err)
		} else {
			// Write the scrambled line
			err = rewriteFromModulo5ToSelected(writer, rec, isOld)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func rewriteFromLoytecToModulo5(writer *csv.Writer, rec []string, channelNumber int) error {
	err := writer.Write([]string{rec[3], strconv.Itoa(channelNumber), "0", loytecToModulo5DataTypeAS(rec[5]), loytecToModulo5DataTypeFS(rec[8]), loytecToModulo5ScalingA(rec[11]), rec[10], loytecToModulo5ByteOrder(rec[12], rec[13], rec[14]), "0", "0", "0", rec[4], loytecToModulo5FunctionCode(rec[5], rec[16]), rec[6], "1", "0"})
	return err
}

func rewriteFromModulo5ToSelected(writer *csv.Writer, rec []string, isOld bool) error {

	if selectedOutput == "modulo6" {
		// Modulo 6
		if isOld {
			return writer.Write([]string{rec[1], rec[2], translateDataTypeAS(rec[4], rec[3]), rec[5], rec[6], rec[11], rec[12], rec[13], translateDataTypeFS(rec[4]), translateByteOrder(rec[7]), translateWordOrder(rec[7]), "0", rec[15], "1", queryInterval, "0", rec[0]})
		} else {
			return writer.Write([]string{rec[1], rec[2], translateDataTypeAS(rec[4], rec[3]), rec[5], rec[6], rec[10], rec[11], rec[12], translateDataTypeFS(rec[4]), translateByteOrder(rec[7]), translateWordOrder(rec[7]), "0", rec[13], "1", queryInterval, "0", rec[0]})
		}
	} else if selectedOutput == "ecos504" {
		// Ecos 504
		if isOld {
			return writer.Write([]string{rec[1], translateDataTypeAS(rec[4], rec[3]), rec[5], rec[6], rec[11], rec[12], rec[13], translateDataTypeFS(rec[4]), translateByteOrder(rec[7]), translateWordOrder(rec[7]), "0", rec[15], "1", queryInterval, "0", rec[0]})
		} else {
			return writer.Write([]string{rec[1], translateDataTypeAS(rec[4], rec[3]), rec[5], rec[6], rec[10], rec[11], rec[12], translateDataTypeFS(rec[4]), translateByteOrder(rec[7]), translateWordOrder(rec[7]), "0", rec[13], "1", queryInterval, "0", rec[0]})
		}
	} else if selectedOutput == "modulo5" {
		// Modulo 5
		if isOld {
			return writer.Write([]string{rec[0], rec[1], rec[2], rec[3], rec[4], rec[5], rec[6], rec[7], rec[8], rec[9], rec[11], rec[12], rec[13], rec[15]})
		} else {
			// If it's already a new version there is no reason to change it
			return writer.Write(rec)
		}
	} else if selectedOutput == "modulo5_old" {
		// Modulo 5 old (v < 1.16)
		if isOld {
			// If it's already an old version there is no reason to change it
			return writer.Write(rec)
		} else {
			return writer.Write([]string{rec[0], rec[1], rec[2], rec[3], rec[4], rec[5], rec[6], rec[7], rec[8], rec[9], "0", rec[10], rec[11], rec[12], "1", rec[13]})
		}
	} else if selectedOutput == "linx" {
		// Modulo 5 old (v < 1.16)
		if isOld {
			swap16, swap32, swap64 := modulo5ToLoytecSwap(rec[7])
			return writer.Write([]string{"UID_Seq", "Modbus Port RS485.Datapoints." + rec[0], rec[0], rec[0], rec[11], modulo5ToLoytecRegisterType(rec[12]), rec[13], modulo5ToLoytecRegisterLength(rec[4]), modulo5ToLoytecModbusDataType(rec[4]), "1", rec[6], modulo5ToLoytecExponent(rec[5]), swap16, swap32, swap64, "", "Input"})
		} else {
			swap16, swap32, swap64 := modulo5ToLoytecSwap(rec[7])
			return writer.Write([]string{"UID_Seq", "Modbus Port RS485.Datapoints." + rec[0], rec[0], rec[0], rec[10], modulo5ToLoytecRegisterType(rec[11]), rec[12], modulo5ToLoytecRegisterLength(rec[4]), modulo5ToLoytecModbusDataType(rec[4]), "1", rec[6], modulo5ToLoytecExponent(rec[5]), swap16, swap32, swap64, "", "Input"})
		}
	} else {
		return errors.New("impossibile riconoscere il tipo di output selezionato")
	}
}

func rewriteFromModulo6ToModulo5Old(writer *csv.Writer, rec []string) error {

	return writer.Write([]string{rec[16], rec[0], rec[1], rec[2], fromModulo6DataTypeFSToModulo5DataTypeFS(rec[6]), rec[3], rec[4], Modulo6ToModulo5ByteOrder(rec[9], rec[10], rec[11]), "0", "0", "0", rec[5], rec[6], rec[7], "1", "0"})

}

func rewriteFromEcos504ToModulo5Old(writer *csv.Writer, rec []string) error {

	return writer.Write([]string{rec[15], rec[0], "0", rec[1], fromModulo6DataTypeFSToModulo5DataTypeFS(rec[5]), rec[2], rec[3], Modulo6ToModulo5ByteOrder(rec[8], rec[9], rec[10]), "0", "0", "0", rec[4], rec[5], rec[6], "1", "0"})

}

func transformModulo6ToModulo5(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) (string, error) {
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
		} else if line == 3 {
			// Write the header of Modulo 5
			err = writer.Write(table_headers["modulo5_old"])
			checkError("Cannot write to file", err)
		} else {
			// Write the scrambled line
			err = rewriteFromModulo6ToModulo5Old(writer, rec)
			checkError("Cannot write to file", err)
			channelNumber++
		}
		if err != nil {
			return "", err
		}
	}

	return outputfile + ".temp", nil

}

func transformEcos504ToModulo5(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) (string, error) {
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
		} else if line == 3 {
			// Write the header of Modulo 5
			err = writer.Write(table_headers["modulo5_old"])
			checkError("Cannot write to file", err)
		} else {
			// Write the scrambled line
			err = rewriteFromEcos504ToModulo5Old(writer, rec)
			checkError("Cannot write to file", err)
			channelNumber++
		}
		if err != nil {
			return "", err
		}
	}

	return outputfile + ".temp", nil

}
