package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

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

	labelInput := widgets.NewQLabel2("<Non selezionato> (Seleziona Modulo 5 o LINX csv)", nil, core.Qt__BypassWindowManagerHint)
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
			isLoytech, err := isLoytech(fileName)

			if err != nil {
				widgets.QMessageBox_Information(widget, "Errore", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			} else {
				if isLoytech {
					fileName, err = transformLoytecToModulo5(fileName, outputfileName, queryInterval, widget)
					if err != nil {
						widgets.QMessageBox_Information(widget, "Errore", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
					}
				}

				// Trasform the input file in the output file
				err = transformFromModulo5ToSelected(fileName, outputfileName, queryInterval, widget)
				if err != nil {
					widgets.QMessageBox_Information(widget, "Errore", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				} else {
					widgets.QMessageBox_Information(widget, "Successo!", "File generato con successo", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
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

func checkError(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
	}
}
