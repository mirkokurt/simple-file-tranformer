package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var table_headers map[string][]string
var selectedOutput string
var queryInterval string
var outputfileName string

func main() {

	var fileName string
	outputfileName = ""

	table_headers = make(map[string][]string)

	table_headers["ecos504"] = []string{"Channel number", "Data type AS [Vas]", "Scaling[A]", "Offset[B]", "Modbus Slave Address", "Data point type", "Start Address", "Modbus data type", "Byte order", "Word order", "DWord order", "Bit selection", "Bit quantity", "Query Interval", "SingleTg", "Description"}
	table_headers["modulo5"] = []string{"Description", "Channel number", "Communication direction", "Data type AS [Vas]", "Data type FS [Vfs]", "Scaling parameter [A]", "Scaling parameter [B]", "Byte order", "Threshold", "Send-Priority", "Modbus Slave address", "Function-Code", "Address", "Bit number"}
	table_headers["modulo6"] = []string{"Channel number", "Communication direction", "Data type AS [Vas]", "Scaling[A]", "Offset[B]", "Modbus Slave Address", "Data point type", "Start Address", "Modbus data type", "Byte order", "Word order", "DWord order", "Bit selection", "Bit quantity", "Query Interval", "SingleTg", "Description"}
	table_headers["modulo5_old"] = []string{"description", "channel nbr", "direction", "DP type AS", "DP type FS", "scaling parameter A", "scaling parameter B", "byte order", "COV inc", "priority", "triggerd", "Slave Adress", "FC", "Adresse", "Anzahl", "BirNr"}

	//set default output
	selectedOutput = "ecos504"

	widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, core.Qt__WindowStaysOnTopHint)
	window.SetWindowTitle("Modulo5 to Modulo6/Ecos504 file converter")
	window.SetMinimumSize2(600, 400)

	logo := widgets.NewQLabel(nil, core.Qt__CoverWindow)
	logo.SetPixmap(gui.NewQPixmap3("logo.png", "", core.Qt__NoAlpha))

	horizontalLayout := widgets.NewQBoxLayout(widgets.QBoxLayout__TopToBottom, nil)
	horizontalLayout.AddWidget(logo, 0, core.Qt__AlignHCenter)

	formLayout := widgets.NewQFormLayout(nil)
	formLayout.SetLabelAlignment(core.Qt__AlignHCenter)

	widgetInside := widgets.NewQWidget(nil, 0)
	widgetInside.SetLayout(formLayout)
	widgetInside.SetMinimumWidth(600)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(horizontalLayout)

	horizontalLayout.AddWidget(widgetInside, 0, core.Qt__AlignHCenter)

	size := core.NewQSize()
	size.SetHeight(20)
	size.SetWidth(150)

	labelInput := widgets.NewQLabel2("<Non selezionato> (Seleziona Modulo 5 csv)", nil, core.Qt__BypassWindowManagerHint)
	labelInput.SetAlignment(core.Qt__AlignHCenter)

	labelQueryInterval := widgets.NewQLabel2("Query interval", nil, core.Qt__BypassWindowManagerHint)
	inputQueryInterval := widgets.NewQLineEdit(nil)
	inputQueryInterval.SetText("240")
	inputQueryInterval.SetAlignment(core.Qt__AlignHCenter)
	inputQueryInterval.SetMinimumSize(size)

	labelOutput := widgets.NewQLabel2("<Non selezionato>", nil, core.Qt__BypassWindowManagerHint)
	labelOutput.SetAlignment(core.Qt__AlignHCenter)

	selectInput := widgets.NewQPushButton2("Scegli file input", nil)
	selectInput.ConnectClicked(func(checked bool) {

		var fileDialog = widgets.NewQFileDialog2(nil, "Open File...", "", "")
		fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
		fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
		var mimeTypes = []string{"text/plain"}
		fileDialog.SetMimeTypeFilters(mimeTypes)

		fileName = fileDialog.GetOpenFileName(nil, "Open File", "C:\\Users\\curtomir\\Desktop", "Comma Separated (*.csv)", "", widgets.QFileDialog__Option(fileDialog.AcceptMode()))
		fmt.Printf("file path is %s", fileName)
		labelInput.SetText(fileName)
		outputfileName = strings.ReplaceAll(fileName, ".csv", "") + "_" + selectedOutput + ".csv"
		labelOutput.SetText(outputfileName)

	})
	selectInput.SetMinimumSize(size)

	selectOutput := widgets.NewQPushButton2("Scegli file output", nil)
	selectOutput.ConnectClicked(func(checked bool) {

		var fileDialog = widgets.NewQFileDialog2(nil, "Open File...", "", "")
		fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
		fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
		var mimeTypes = []string{"text/plain"}
		fileDialog.SetMimeTypeFilters(mimeTypes)

		outputfileName = fileDialog.GetSaveFileName(nil, "Save File", "C:\\Users\\curtomir\\Desktop", "Comma Separated (*.csv)", "", widgets.QFileDialog__ShowDirsOnly)
		fmt.Printf("file path is %s", outputfileName)
		labelOutput.SetText(outputfileName)
	})
	selectOutput.SetMinimumSize(size)

	tranformlabel := widgets.NewQLabel2("", nil, core.Qt__BypassWindowManagerHint)
	tranformButton := widgets.NewQPushButton2("Trasforma", nil)
	tranformButton.ConnectClicked(func(checked bool) {
		queryInterval = inputQueryInterval.Text()
		_, err := strconv.Atoi(queryInterval)
		if err != nil {
			widgets.QMessageBox_Warning(widget, "Errore", "Inserisci un Query Interval valido", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else if fileName == "" || outputfileName == "" {
			widgets.QMessageBox_Warning(widget, "Errore", "Seleziona prima i file!", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			// Check if it is a Loytec input file first and in the case perform a pre-transformation in Modulo5 file format
			fileName, err = checkLoytec(fileName, outputfileName, queryInterval, widget)
			if err != nil {
				widgets.QMessageBox_Information(widget, "Errore", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			}

			// Trasform the input file in the output file
			err := transform(fileName, outputfileName, queryInterval, widget)
			if err != nil {
				widgets.QMessageBox_Information(widget, "Errore", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			} else {
				widgets.QMessageBox_Information(widget, "Successo!", "File generato con successo", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			}
		}
	})

	labelCombo := widgets.NewQLabel2("Seleziona il tipo di output", nil, core.Qt__BypassWindowManagerHint)
	labelCombo.SetAlignment(core.Qt__AlignHCenter)

	comboOutput := widgets.NewQComboBox(nil)
	comboOutput.SetEnabled(true)
	comboOutput.SetEditable(true)
	comboOutput.LineEdit().SetReadOnly(true)
	comboOutput.LineEdit().SetAlignment(core.Qt__AlignHCenter)
	comboOutput.AddItems([]string{"ecos504", "modulo6", "modulo5", "modulo5_old"})
	comboOutput.ConnectCurrentIndexChanged(func(index int) { setOutputType(index, labelOutput) })

	formLayout.AddRow(selectInput, labelInput)
	formLayout.AddRow(selectOutput, labelOutput)
	formLayout.AddRow(inputQueryInterval, labelQueryInterval)
	formLayout.AddRow(labelCombo, comboOutput)
	formLayout.AddRow(tranformlabel, tranformButton)

	window.SetCentralWidget(widget)
	window.Show()

	widgets.QApplication_Exec()
}

func setOutputType(index int, labelOutput *widgets.QLabel) {

	switch index {
	case 0:
		selectedOutput = "ecos504"
	case 1:
		selectedOutput = "modulo6"
	case 2:
		selectedOutput = "modulo5"
	case 3:
		selectedOutput = "modulo5_old"
	default:
		selectedOutput = "ecos504"
	}
	if outputfileName != "" {
		outputfileName = strings.ReplaceAll(outputfileName, "_ecos504", "")
		outputfileName = strings.ReplaceAll(outputfileName, "_modulo6", "")
		outputfileName = strings.ReplaceAll(outputfileName, "_modulo5_old", "")
		outputfileName = strings.ReplaceAll(outputfileName, "_modulo5", "")
		outputfileName = strings.ReplaceAll(outputfileName, ".csv", "") + "_" + selectedOutput + ".csv"
		labelOutput.SetText(outputfileName)
	}

}

func checkLoytec(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) (string, error) {
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

	rec, err := csvReader.Read()
	if err == io.EOF || !strings.Contains(rec[0], "LOYTEC") {
		// If the input file is not a LOYTEC, return the input file and start a normal conversion
		return inputfile, nil
	} else {
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
		line := 1
		channelNumber := 0
		for {
			// Increment the line
			line++
			rec, err = csvReader.Read()
			if err == io.EOF {
				break
			}
			checkError("Cannot read the file", err)

			// If one of the first two lines
			if line < 4 {
				empty_rec := []string{"", "", ""}
				err = writer.Write(empty_rec)
				checkError("Cannot write to file", err)
			} else if line >= 4 && line < 7 {
				// Skip the lines
			} else if line == 7 {
				// Write the header of Modulo 5
				err = writer.Write(table_headers["modulo5_old"])
				checkError("Cannot write to file", err)
			} else {
				// Write the scrambled line
				err = rewriteContentFoRLoytec(writer, rec, channelNumber)
				checkError("Cannot write to file", err)
				channelNumber++
			}
			if err != nil {
				return "", err
			}
		}

		return outputfile + ".temp", nil
	}
}

func transform(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) error {

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
			isOld = detectTypeFromHeader(rec[10])

			// Write the selected device table header
			err = writer.Write(table_headers[selectedOutput])
			checkError("Cannot write to file", err)
		} else {
			// Write the scrambled line
			err = rewriteContent(writer, rec, isOld)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func rewriteContentFoRLoytec(writer *csv.Writer, rec []string, channelNumber int) error {
	err := writer.Write([]string{rec[3], strconv.Itoa(channelNumber), "0", loytecToModulo5DataTypeAS(rec[5]), loytecToModulo5DataTypeFS(rec[8]), rec[10], rec[9], "0", "0", "0", "0", rec[4], loytecToModulo5FunctionCode(rec[5], rec[16]), rec[6], "1", "0"})
	return err
}

func rewriteContent(writer *csv.Writer, rec []string, isOld bool) error {

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
	} else {
		return errors.New("impossibile riconoscere il tipo di output selezionato")
	}
}

func detectTypeFromHeader(check string) bool {
	if check == "triggerd" || check == "Trigger: mode" {
		return true
	}
	return false
}

func checkError(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
	}
}

func translateDataTypeFS(input string) string {
	var output string
	switch input {
	case "0":
		output = "0"
	case "1":
		output = "1"
	case "2":
		output = "5"
	case "3":
		output = "2"
	case "4":
		output = "6"
	case "5":
		output = "3"
	case "6":
		output = "7"
	case "7":
		output = "1"
	case "8":
		output = "3"
	case "9":
		output = "5"
	case "10":
		output = "9"
	case "11":
		output = "TEXT"
	case "12":
		output = "4"
	case "13":
		output = "8"
	case "14":
		output = "10"
	default:
		output = "NaN"
	}

	return output
}

func translateDataTypeAS(dataTypeFS string, DdataTypeAS string) string {
	var output string
	switch dataTypeFS {
	case "7":
		output = "2"
	case "8":
		output = "2"
	case "9":
		output = "2"
	default:
		output = DdataTypeAS
	}

	return output
}

func loytecToModulo5DataTypeAS(registerType string) string {
	var output string
	switch registerType {
	case "INPUT":
		output = "2"
	case "DISCRETE":
		output = "2"
	case "HOLD":
		output = "0"
	default:
		output = "0"
	}

	return output
}

func loytecToModulo5DataTypeFS(dataType string) string {
	var output string
	switch dataType {
	case "uint8":
		output = "0"
	case "uint16":
		output = "1"
	case "uint32":
		output = "2"
	case "uint64":
		output = "3"
	case "int8":
		output = "4"
	case "int16":
		output = "5"
	case "int32":
		output = "6"
	case "int64":
		output = "7"
	case "float32":
		output = "8"
	case "float64":
		output = "9"
	default:
		output = "0"
	}

	return output
}

func loytecToModulo5FunctionCode(registerType string, Direction string) string {
	output := "0"

	if registerType == "DISCRETE" {
		output = "2"
	}

	if registerType == "COIL" {
		if Direction == "Input" {
			output = "1"
		} else if Direction == "Output" {
			output = "5"
		} else if Direction == "Value" {
			output = "5"
		}
	}

	if registerType == "INPUT" {
		output = "4"
	}

	if registerType == "HOLD" {
		if Direction == "Input" {
			output = "3"
		} else if Direction == "Output" {
			output = "6"
		} else if Direction == "Value" {
			output = "6"
		}
	}

	return output
}

// Just a mock translation for future implementation
func translateByteOrder(byteOrder string) string {
	return "1"
}

// Just a mock translation for future implementation
func translateWordOrder(byteOrder string) string {
	return byteOrder
}
