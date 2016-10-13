package edf

// #include "C/oa.h"
// #include "C/buffer.h"
// #include "C/csv2ascii.h"
import "C"
import "os"
// import "log"
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
