package edf

// #include "C/oa.h"
// #include "C/buffer.h"
// #include "C/csv2ascii.h"
import "C"

// TODO Make this repository not depend upon github.com/ishiikurisu/OA
func Csv2SingleWoC(inlet string) {

}

func Csv2MultipleWoC(inlet string) {
	
}

// Converts the generated *.csv file from WriteCSV to a single *.ascii file
// with every channel recording.
func Csv2Single(inlet string) {
    C.csv2single(C.CString(inlet))
}

// Converts the generated *.csv file from WriteCSV to multiple *.ascii files,
// as many channels exist in the recording.
func Csv2Multiple(inlet string) {
    C.csv2multiple(C.CString(inlet))
}
