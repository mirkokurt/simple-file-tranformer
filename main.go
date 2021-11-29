package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

type fileInfo struct {
	complete_path string
	name          string
}

var rule_sets []rule

var (
	dir_path        = flag.String("dir_path", "/", "This is the path of the directory that contains the file you want to ")
	ext_type        = flag.String("ext_type", ".csv", "This is the extension of the file to search for")
	output_ext_type = flag.String("output_ext_type", ".gr", "This is the extension of the file to search for")
)

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
		//checkError("Cannot read the file", err)

		// Wait one second before checking for the files again
		time.Sleep(1000 * time.Millisecond)
	}

}

func processFile(toBeProcessed chan fileInfo) error {
	for {
		//Wait for a new file found
		file_tpb := <-toBeProcessed

		fmt.Printf("Processing the file %s \n", file_tpb.name)

		// Open the file in input
		fi, err := os.Open(file_tpb.complete_path)
		if err != nil {
			fmt.Printf("An error occurred opening the file: %s \n", err)
			continue
		}

		// Create the path of the new file
		new_file_path := strings.ReplaceAll(file_tpb.complete_path, *ext_type, *output_ext_type)

		fmt.Printf("Creating the file %s \n", new_file_path)
		// Open or create the output file
		fo, err := os.Create(new_file_path)
		if err != nil {
			fmt.Printf("An error occurred while creating the output file: %s \n", err)
			continue
		}
		// Close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				fmt.Printf("An error occurred while closing the file: %s \n", err)
				return
			}
		}()

		// Create a writer
		writer := bufio.NewWriter(fo)

		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			// Take a line from the file
			line := scanner.Text()
			writer.WriteString(line + "\n")
		}
		writer.Flush()

		if err := scanner.Err(); err != nil {
			fmt.Printf("An error occurred scanning the input file, the file will not be deleted: %s \n", err)
			continue
		}

		fi.Close()
		// Removing file from the directory
		// Using Remove() function
		err = os.Remove(file_tpb.complete_path)
		if err != nil {
			fmt.Printf("The file %s cannot be removed due to error: %s \n", file_tpb.name, err)
		}

	}
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
