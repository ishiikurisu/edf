package edf

// #include "C/oa.h"
// #include "C/buffer.h"
// #include "C/csv2ascii.h"
import "C"
import "os"
import "log"
import "strings"
import "bufio"

// TODO Make this repository not depend upon github.com/ishiikurisu/OA

// Converts the given file from the CSV format, as produced by edf.WriteCSV to
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

// TODO Write this monster
func Csv2Multiple(inlet string) {
	inputFile, _ := os.Open(inlet)
	scanner := bufio.NewScanner(inputFile)
	defer inputFile.Close()

	// Extracting data from header
	scanner.Scan()
	header := scanner.Text()
	labels := extractLabelsFromHeader(header)

	// Opening output buffers
	// TODO Generate output names
	// TODO Open files and defer their closing

	// Writing data to each channel
	// TODO While scanning input, write each information in their respetive
	// TODO Flush output buffers

	log.Printf("%#v\n", labels)
	C.csv2multiple(C.CString(inlet))
}

// Converts the generated *.csv file from WriteCSV to a single *.ascii file
// with every channel recording.
func Csv2SingleWithC(inlet string) {
    C.csv2single(C.CString(inlet))
}

// Converts the generated *.csv file from WriteCSV to multiple *.ascii files,
// as many channels exist in the recording.
func Csv2MultipleWithC(inlet string) {
    C.csv2multiple(C.CString(inlet))
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
		} else if (stuff[0] == "chan") {
			noChans = str2int(stuff[1])
		}
	}

	// Extracting labels
	fields = separateString(labels[1], noChans)
	labels = make([]string, 0)
	for _, field := range fields {
		field = trimField(field)
		if (field != "EDF Annotations") {
			labels = addItem(labels, field) // I CAN'T BELIEVE APPEND DOES NOT WORK HERE
		}
	}

	return labels
}

func trimField(inlet string) string {
	lower := 0
	upper := len(inlet) - 1

	for inlet[lower] == ' ' || inlet[lower] == '[' {
		lower++
	}
	for inlet[upper] == ' ' {
		upper--
	}

	return inlet[lower:(upper+1)]
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
