package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type fileInfo struct {
	complete_path string
	name          string
}

type Record struct {
	Codice_soggetto  string `json:"cod. soggetto"`
	Tipo_Documento   string `json:"tipo documento"`
	Numero_Documento string `json:"numero documento"`
	Data_Documento   string `json:"data documento"`
	Descr_parte1     string `json:"descrizione documento (parte1)"`
	Descr_parte2     string `json:"descrizione documento (parte2)"`
	Descr_parte3     string `json:"descrizione documento (parte3)"`
	Descr_parte4     string `json:"descrizione documento (parte4)"`
	Descr_parte5     string `json:"descrizione documento (parte5)"`
	Descr_parte6     string `json:"descrizione documento (parte6)"`
	Descr_parte7     string `json:"descrizione documento (parte7)"`
	Descr_parte8     string `json:"descrizione documento (parte8)"`
	Descr_parte9     string `json:"descrizione documento (parte9)"`
	Descr_parte10    string `json:"descrizione documento (parte10)"`
}

var (
	dir_path        = flag.String("dir_path", "/", "This is the path of the directory that contains the file you want to ")
	ext_type        = flag.String("ext_type", ".json", "This is the extension of the file to search for")
	output_ext_type = flag.String("output_ext_type", ".gr", "This is the extension of the file to search for")
)

const Tab string = "\t"

var globalCounter = 0

func main() {

	// Take parameters from the command line

	flag.Parse()

	if len(os.Args) < 2 || *dir_path == "" || *ext_type == "" {
		usage("Please don't insert empty arguments!!")
	}

	toBeProcessed := make(chan fileInfo, 1)
	go processFile(toBeProcessed)

	for {
		// Check in the path for a file with the defined type
		err := filepath.Walk(*dir_path, func(path string, info os.FileInfo, err error) error {
			// If the file has the correct extension
			fileExtension := filepath.Ext(path)
			if fileExtension == *ext_type {
				toBeProcessed <- fileInfo{path, info.Name()}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

		// Wait one second before checking for the files again
		time.Sleep(10000 * time.Millisecond)
	}

}

func processFile(toBeProcessed chan fileInfo) error {
	for {
		//Wait for a new file found
		file_tpb := <-toBeProcessed

		fmt.Printf("Processing the file %s \n", file_tpb.name)

		file, _ := ioutil.ReadFile(file_tpb.complete_path)

		// Create an empty record data
		rd := Record{}

		//Marshal the file content into the empty json structure
		err := json.Unmarshal([]byte(file), &rd)
		if err != nil {
			fmt.Printf("The file %s cannot be processed due to error: %s \n", file_tpb.name, err)
		}

		globalCounter++
		// Create the path of the new file
		path := strings.ReplaceAll(file_tpb.complete_path, file_tpb.name, "")
		new_file_path := fmt.Sprintf("%s%s%s%d%s", path, "gen", rd.Tipo_Documento, globalCounter, *output_ext_type)

		fmt.Printf("Creating the file %s \n", new_file_path)
		// Open or create the output file
		fo, err := os.Create(new_file_path)
		if err != nil {
			fmt.Printf("An error occurred while creating the output file: %s \n", err)
			continue
		}

		// Create a writer
		writer := bufio.NewWriter(fo)

		cleanDocumentNumber := strings.ReplaceAll(rd.Numero_Documento, "-", "")

		writer.WriteString("I" + Tab +
			"034" + Tab +
			shortDocType(rd.Tipo_Documento) + Tab +
			strings.ReplaceAll(rd.Codice_soggetto, "34-", "") + Tab +
			rd.Tipo_Documento + Tab +
			cleanDocumentNumber + Tab +
			"*" + Tab +
			formatDate(rd.Data_Documento) + Tab +
			// Manage differences between DDTC and others
			inline_if(rd.Tipo_Documento == "DDTC", rd.Descr_parte1, rd.Descr_parte1).(string) + Tab +
			printDescritpion(rd, rd.Tipo_Documento) + Tab +
			strings.ReplaceAll(file_tpb.name, *ext_type, ".pdf") + Tab +
			"0" + Tab +
			"*" + Tab +
			"0" + Tab +
			"0" + Tab +
			"*" + Tab +
			"aida" + Tab +
			"X" +
			"\r\n")

		writer.Flush()

		// TODO: in the production version the file must be removed
		// Removing file from the directory
		/*err = os.Remove(file_tpb.complete_path)
		if err != nil {
			fmt.Printf("The file %s cannot be removed due to error: %s \n", file_tpb.name, err)
		}*/

		// Close fo on exit and check for its returned error
		if err := fo.Close(); err != nil {
			fmt.Printf("An error occurred while closing the file: %s \n", err)
			continue
		}

		// Rename the file
		err = os.Rename(file_tpb.complete_path, strings.ReplaceAll(file_tpb.complete_path, *ext_type, ".bck"))
		if err != nil {
			fmt.Printf("Error renaming the file: %s \n", err)
			continue
		}

	}
}

//inline version of if_then_else
func inline_if(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func shortDocType(docType string) string {
	switch docType {
	case "DDTF":
		return "F"
	case "DDTC":
		return "C"
	default:
		return "F"
	}
}

func formatDate(date string) string {
	i, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		fmt.Println("Error formatting the date")
		return ""
	}
	tm := time.UnixMilli(i)
	//fmt.Println(tm)
	return fmt.Sprintf(tm.Format("02/01/2006"))
}

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <dir_path> <ext_type> [output_ext_type]\n"+
			"       where <dir_path> is the directory where search the input files \n"+
			"       and <ext_type> is the extension of the files to search for\n"+
			"       and <output_ext_type> is the exntesion of the output files to be created\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func printDescritpion(rd Record, docType string) string {
	description := ""

	fields := reflect.TypeOf(rd)
	values := reflect.ValueOf(rd)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i).Interface().(string)
		if strings.HasPrefix(field.Name, "Descr_") && value != "" {
			// If the doc type is "DDTC" skip the first description field
			if docType == "DDTC" && field.Name == "Descr_parte1" {
				continue
			}
			if len(description) == 0 {
				description += value
			} else {
				description += " - " + value
			}
		}
	}
	return description
}
