package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type rule struct {
	re           string
	input_string string
	replace_with string
}

type table_header struct {
	dev_type string
	header   []string
}

var rule_sets []rule
var table_headers map[string][]string
var selectedOutput string

func main() {

	var fileName string
	var outputfileName string

	table_headers = make(map[string][]string)

	table_headers["ecos504"] = []string{"Channel number", "Data type AS [Vas]", "Scaling[A]", "Offset[B]", "Modbus Slave Address", "Data point type", "Start Address", "Modbus data type", "Byte order", "Word order", "DWord order", "Bit selection", "Bit quantity", "Query Interval", "SingleTg", "Description"}
	table_headers["modulo5"] = []string{"Description", "Channel number", "Communication direction", "Data type AS [Vas]", "Data type FS [Vfs]", "Scaling parameter [A]", "Scaling parameter [B]", "Byte order", "Threshold", "Send-Priority", "Modbus Slave address", "Function-Code", "Address", "Bit number"}
	table_headers["modulo6"] = []string{"Channel number", "Communication direction", "Data type AS [Vas]", "Scaling[A]", "Offset[B]", "Modbus Slave Address", "Data point type", "Start Address", "Modbus data type", "Byte order", "Word order", "DWord order", "Bit selection", "Bit quantity", "Query Interval", "SingleTg", "Description"}

	//set default output
	selectedOutput = "modulo6"

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

	})
	selectInput.SetMinimumSize(size)

	labelOutput := widgets.NewQLabel2("<Non selezionato>", nil, core.Qt__BypassWindowManagerHint)
	labelOutput.SetAlignment(core.Qt__AlignHCenter)

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

	//box := widgets.NewQMessageBox(nil)

	tranformlabel := widgets.NewQLabel2("", nil, core.Qt__BypassWindowManagerHint)
	tranformButton := widgets.NewQPushButton2("Trasforma", nil)
	tranformButton.ConnectClicked(func(checked bool) {
		queryInterval := inputQueryInterval.Text()
		_, err := strconv.Atoi(queryInterval)
		if err != nil {
			widgets.QMessageBox_Warning(widget, "Errore", "Inserisci un Query Interval valido", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else if fileName == "" || outputfileName == "" {
			widgets.QMessageBox_Warning(widget, "Errore", "Seleziona prima i file!", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			transform(fileName, outputfileName, queryInterval, widget)
			widgets.QMessageBox_Information(widget, "Errore", "File generato con successo", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})

	labelCombo := widgets.NewQLabel2("Seleziona il tipo di output", nil, core.Qt__BypassWindowManagerHint)
	labelCombo.SetAlignment(core.Qt__AlignHCenter)

	comboOutput := widgets.NewQComboBox(nil)
	comboOutput.SetEnabled(true)
	comboOutput.SetEditable(true)
	comboOutput.LineEdit().SetReadOnly(true)
	comboOutput.LineEdit().SetAlignment(core.Qt__AlignHCenter)
	comboOutput.AddItems([]string{"ecos504", "modulo6"})
	comboOutput.ConnectCurrentIndexChanged(func(index int) { setOutputType(index) })

	formLayout.AddRow(selectInput, labelInput)
	formLayout.AddRow(selectOutput, labelOutput)
	formLayout.AddRow(inputQueryInterval, labelQueryInterval)
	formLayout.AddRow(labelCombo, comboOutput)
	formLayout.AddRow(tranformlabel, tranformButton)
	//formLayout.AddRow(labelCombo, logo)

	window.SetCentralWidget(widget)
	window.Show()

	widgets.QApplication_Exec()
}

func setOutputType(index int) {

	switch index {
	case 0:
		selectedOutput = "ecos504"
	case 1:
		selectedOutput = "modulo6"
	default:
		selectedOutput = "ecos504"
	}

}

func transform(inputfile string, outputfile string, queryInterval string, widget *widgets.QWidget) {

	// Open the input file
	fr, err := os.Open(inputfile)
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer fr.Close()

	// Open the output file
	fo, err := os.Create(outputfile)
	if err != nil {
		widgets.QMessageBox_Information(widget, "Errore", "Attenzione il file di output è aperto e non può essere scritto", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	// Close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a writer
	writer := csv.NewWriter(fo)
	defer writer.Flush()

	writer.Comma = ';'

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(fr)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'

	// Init the line counter
	line := 0

	for {
		// Increment the line
		line++
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		checkError("Cannot read the file", err)

		// Write the output file's line

		// If one of the first tweo lines
		if line < 3 {
			// Write the exact same line
			err = writer.Write(rec)
			checkError("Cannot write to file", err)
		} else if line == 3 {
			// Write the selected device table header
			err = writer.Write(table_headers[selectedOutput])
			checkError("Cannot write to file", err)
		} else {
			// Write the scrambled line
			if selectedOutput == "modulo6" {
				// Modulo 6
				err = writer.Write([]string{rec[1], rec[2], translateDataTypeAS(rec[4], rec[3]), rec[5], rec[6], rec[10], rec[11], rec[12], translateDataTypeFS(rec[4]), rec[7], rec[7], "0", rec[13], "1", queryInterval, "0", rec[0]})
			} else if selectedOutput == "ecos504" {
				// Ecos 504
				err = writer.Write([]string{rec[1], translateDataTypeAS(rec[4], rec[3]), rec[5], rec[6], rec[10], rec[11], rec[12], translateDataTypeFS(rec[4]), rec[7], rec[7], "0", rec[13], "1", queryInterval, "0", rec[0]})
				checkError("Cannot write to file", err)
			} else {
				widgets.QMessageBox_Warning(widget, "Errore", "Seleziona la tipologia di file di output", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			}

		}
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func translateDataTypeFS(input string) string {
	var output string
	switch input {
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
