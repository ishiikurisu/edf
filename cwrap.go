package edf

// #include "C/oa.h"
// #include "C/buffer.h"
// #include "C/csv2ascii.h"
import "C"
// import "os"
import "log"
import "strings"
// import "bufio"

// TODO Make this repository not depend upon github.com/ishiikurisu/OA
func Csv2Single(inlet string) {
	outlet := generateSingleOutput(inlet)
	// inputFile, _ := os.Open(inlet)
	// outputFile, _ := os.OpenFile(outlet)
	// scanner := bufio.NewScanner(fp)
	// typewriter := bufio.NewWriter(os.OpenFile(outlet, ))
	// defer fp.Close()

	log.Println(outlet)
	Csv2SingleWC(inlet)
}

func Csv2Multiple(inlet string) {
	C.csv2multiple(C.CString(inlet))
}

// Converts the generated *.csv file from WriteCSV to a single *.ascii file
// with every channel recording.
func Csv2SingleWC(inlet string) {
    C.csv2single(C.CString(inlet))
}

// Converts the generated *.csv file from WriteCSV to multiple *.ascii files,
// as many channels exist in the recording.
func Csv2MultipleWC(inlet string) {
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

func getLabelsFromHeader(header string)[]string {
	fields := strings.Split(header, ";")

	// Getting labels field

	// Getting number channels

	// Separating labels' names

	return fields
}

func getNumberChan(header string) int {
	return 0
}
