package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func main() {

	var fileName string
	var outputfileName string

	widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("File transformer")
	window.SetMinimumSize2(600, 400)

	formLayout := widgets.NewQFormLayout(nil)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(formLayout)

	input := widgets.NewQLabel2("<Non selezionato>", nil, core.Qt__BypassWindowManagerHint)

	button := widgets.NewQPushButton2("Seleziona", nil)
	button.ConnectClicked(func(checked bool) {
		//widgets.QMessageBox_Information(nil, "Titolo del messaggio", input.Text(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)

		var fileDialog = widgets.NewQFileDialog2(nil, "Open File...", "", "")
		fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
		fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
		var mimeTypes = []string{"text/plain"}
		fileDialog.SetMimeTypeFilters(mimeTypes)

		fileName = fileDialog.GetOpenFileName(nil, "Open File", "C:\\Users\\curtomir\\Desktop", "Text/plain (*.txt)", "", widgets.QFileDialog__Option(fileDialog.AcceptMode()))
		fmt.Printf("file path is %s", fileName)
		input.SetText(fileName)

	})

	output := widgets.NewQLabel2("<Non selezionato>", nil, core.Qt__BypassWindowManagerHint)

	button2 := widgets.NewQPushButton2("Scegli file output", nil)
	button2.ConnectClicked(func(checked bool) {
		//widgets.QMessageBox_Information(nil, "Titolo del messaggio", input.Text(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)

		var fileDialog = widgets.NewQFileDialog2(nil, "Open File...", "", "")
		fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
		fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
		var mimeTypes = []string{"text/plain"}
		fileDialog.SetMimeTypeFilters(mimeTypes)

		outputfileName = fileDialog.GetSaveFileName(nil, "Save File", "C:\\Users\\curtomir\\Desktop", "Text/plain (*.txt)", "", widgets.QFileDialog__ShowDirsOnly)
		fmt.Printf("file path is %s", outputfileName)
		output.SetText(outputfileName)
	})

	buttonlabel := widgets.NewQLabel2("", nil, core.Qt__BypassWindowManagerHint)
	button3 := widgets.NewQPushButton2("Trasforma", nil)
	button3.ConnectClicked(func(checked bool) {
		if fileName == "" || outputfileName == "" {
			widgets.QMessageBox_Information(nil, "Errore", "Seleziona prima i file!", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			transform(fileName, outputfileName)
		}
	})

	formLayout.AddRow(input, button)
	formLayout.AddRow(output, button2)
	formLayout.AddRow(buttonlabel, button3)

	window.SetCentralWidget(widget)
	window.Show()

	widgets.QApplication_Exec()
}

type rule struct {
	re           string
	input_string string
	replace_with string
}

var rule_sets []rule

func transform(inputfile string, outputfile string) {

	// Open the file with the rules to apply
	fr, err := os.Open("rules.csv")
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer fr.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(fr)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Init the rule set
		rul := rule{rec[0], rec[1], rec[2]}
		rule_sets = append(rule_sets, rul)
	}

	// Open the input file
	fi, err := os.Open(inputfile)
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer fi.Close()

	// Open the output file
	fo, err := os.Create(outputfile)
	if err != nil {
		panic(err)
	}
	// Close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a writer
	writer := bufio.NewWriter(fo)

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		var res []string
		// Take a line from the file
		line := scanner.Text()
		// Iterate all the rules in the set to search on rule to apply
		for _, elem := range rule_sets {
			res = regexp.MustCompile(elem.re).FindStringSubmatch(line)
			if len(res) > 0 {
				// Apply the replace_with string
				output_string := res[1] + " - " + elem.replace_with
				fmt.Println(output_string)
				writeFile(writer, output_string)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	writer.Flush()
}

func writeFile(writer *bufio.Writer, output_string string) {
	writer.WriteString(output_string + "\n")
}
