package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type rule struct {
	re           string
	input_string string
	replace_with string
}

var rule_sets []rule

func main() {

	//Open the file with the rules to apply
	fr, err := os.Open("rules.csv")
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer fr.Close()

	//Read csv values using csv.Reader
	csvReader := csv.NewReader(fr)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//Init the rule set
		rul := rule{rec[0], rec[1], rec[2]}
		rule_sets = append(rule_sets, rul)
	}

	//Open the input file
	fi, err := os.Open("input.txt")
	if err != nil {
		fmt.Print("There has been an error!: ", err)
	}
	defer fi.Close()

	//Open the output file
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
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
		//take a line from the file
		line := scanner.Text()
		//iterate all the rules in the set to search on rule to apply
		for _, elem := range rule_sets {
			res = regexp.MustCompile(elem.re).FindStringSubmatch(line)
			if len(res) > 0 {
				//apply the replace_with string
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
