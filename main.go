package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type rule struct {
	re           *regexp.Regexp
	input_string string
	replace_with string
}

var regolaA = rule{regexp.MustCompile(`^(\S*) - REQUEST-A$`), "REQUEST-A", "REQUEST-Z"}
var regolaB = rule{regexp.MustCompile(`^(\S*) - REQUEST-B$`), "REQUEST-B", "REQUEST-X"}
var regolaC = rule{regexp.MustCompile(`^(\S*) - REQUEST-C$`), "REQUEST-C", "REQUEST-Y"}

var rule_sets [3]rule

func main() {

	rule_sets[0] = regolaA
	rule_sets[1] = regolaB
	rule_sets[2] = regolaC

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
			res = elem.re.FindStringSubmatch(line)
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
