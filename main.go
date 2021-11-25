package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

func main() {
	var files []string

	// Take parameters from the command line
	path := flag.String("path", "/", "This is the path of the directory that contains the file you want to ")
	ext_type := flag.String("ext_type", ".cg", "This is the extension of the file to search for")
	flag.Parse()

	if len(os.Args) < 2 || *path == "" || *ext_type == "" {
		usage("Please don't insert empty arguments!!")
	}

	// Check in the path for a file with the defined type
	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		// If the file has the correct extension
		fileExtension := filepath.Ext(path)
		if fileExtension == *ext_type {
			files = append(files, path)
			processFile(path, info.Name(), *ext_type)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		// If the file has the correct extension

		// Open the file
		/*fi, err := os.Open(inputfile)
		if err != nil {
			return err
		}
		defer fi.Close()

		// Open the output file
		fo, err := os.Create(outputfile)
		if err != nil {
			fmt.Printf("Error creating the output file: ", err)
		}
		// Close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				fmt.Printf("Error closing output file: ", err)
			}
		}()*/
		fmt.Println(file)
	}
	//checkError("Cannot read the file", err)
}

func processFile(path, name, ext_type string) error {
	// Open the file in input
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	//find the position of the dot
	dot := len(path) - len(ext_type)
	// Open or create the output file
	fo, err := os.Create(path[0:dot] + "_modified" + path[dot:])
	if err != nil {
		return err
	}
	// Close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			return
		}
	}()

	// Create a writer
	writer := bufio.NewWriter(fo)
	defer writer.Flush()

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		//var res []string
		// Take a line from the file
		line := scanner.Text()
		writer.WriteString(line + "\n")
		/*
			// Iterate all the rules in the set to search on rule to apply
			for _, elem := range rule_sets {
				res = regexp.MustCompile(elem.re).FindStringSubmatch(line)
				if len(res) > 0 {
					// Apply the replace_with string
					output_string := res[1] + " - " + elem.replace_with
					fmt.Println(output_string)
					writer.WriteString(output_string + "\n")
				}
			}
		*/
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return nil
}

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command> <password>\n"+
			"       where <username> is the username to be created (not encrypted) \n"+
			"       and <password> is the password you want to encrypt\n",
		errmsg, os.Args[0])
	os.Exit(2)
}
