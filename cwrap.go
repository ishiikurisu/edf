package edf

import "os"
import "fmt"
import "strings"
import "bufio"

// Csv2Single converts the given file from the CSV format, as produced by edf.WriteCSV to
// another, as produced by edf.WriteASCII
func Csv2Single(inlet string) {
	outlet := generateSingleOutput(inlet)
	inputFile, _ := os.Open(inlet)
	outputFile, _ := os.Create(outlet)
	scanner := bufio.NewScanner(inputFile)
	typewriter := bufio.NewWriter(outputFile)
	defer inputFile.Close()
	defer outputFile.Close()

	// Ignoring header
	scanner.Scan()
	scanner.Text()

	// Extracting data
	for scanner.Scan() {
		from := scanner.Text()
		to := strings.Join(strings.Split(from, ";"), " ")
		typewriter.WriteString(to + "\n")
	}

	typewriter.Flush()
}

// Csv2Multiple converts a single CSV file to multiple ASCII files, each one for a different channel
func Csv2Multiple(inlet string) {
	inputFile, _ := os.Open(inlet)
	scanner := bufio.NewScanner(inputFile)
	defer inputFile.Close()

	// Extracting data from header
	scanner.Scan()
	header := scanner.Text()
	labels := extractLabelsFromHeader(header)

	// Opening output buffers
	outlets := generateMultipleOutputs(inlet, labels)
	outputFiles := createOutputFiles(outlets)
	for _, outputFile := range outputFiles {
		defer outputFile.Close()
	}

	// Writing data to each channel
	for scanner.Scan() {
		from := scanner.Text()
		to := strings.Split(from, ";")
		for i, it := range to {
			outputFiles[i].WriteString(trimField(it) + "\n")
		}
	}
}

/* ######################
   # AUXILIAR FUNCTIONS #
   ###################### */

func generateSingleOutput(inlet string) string {
	index := len(inlet) - 1

	for inlet[index] != '.' {
		index--
	}

	outlet := inlet[0:index] + ".ascii"
	return outlet
}

func extractLabelsFromHeader(header string) []string {
	var fields []string
	var labels []string
	var noChans int

	// Getting labels field
	fields = strings.Split(header, ";")
	for i := 0; i < len(fields); i++ {
		stuff := strings.Split(fields[i], ":")
		if stuff[0] == "labels" {
			labels = stuff
		} else if stuff[0] == "chan" {
			noChans = str2int(stuff[1])
		}
	}

	// Extracting labels
	fields = separateString(labels[1], noChans)
	labels = make([]string, 0)
	for _, field := range fields {
		field = trimField(field)
		if field != "EDF Annotations" {
			labels = addItem(labels, field) // I CAN'T BELIEVE APPEND DOES NOT WORK HERE
		}
	}

	return labels
}

func trimField(inlet string) string {
	lower := 0
	upper := len(inlet) - 1

	if len(inlet) == 0 {
		return inlet
	}

	for inlet[lower] == ' ' || inlet[lower] == '[' {
		lower++
	}
	for inlet[upper] == ' ' {
		upper--
	}

	return inlet[lower:(upper + 1)]
}

func addItem(box []string, item string) []string {
	limit := len(box)
	newBox := make([]string, limit+1)

	for i, it := range box {
		newBox[i] = it
	}

	newBox[limit] = item
	return newBox
}

func generateMultipleOutputs(inlet string, labels []string) []string {
	outlets := make([]string, len(labels))
	raw := inlet
	i := 0

	// Getting to the point
	for i = len(inlet) - 1; inlet[i] != '.'; i-- {

	}
	raw = inlet[0:i]

	// Constructing output
	for j, label := range labels {
		outlets[j] = fmt.Sprintf("%s.%s.ascii", raw, label)
	}

	return outlets
}

func createOutputFiles(outlets []string) []*os.File {
	outputs := make([]*os.File, len(outlets))

	for i, outlet := range outlets {
		outputs[i], _ = os.Create(outlet)
	}

	return outputs
}
