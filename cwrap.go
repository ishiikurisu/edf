package edf

// #include "C/oa.h"
// #include "C/buffer.h"
// #include "C/csv2ascii.h"
import "C"

func Csv2Single(inlet string) {
    C.csv2single(C.CString(inlet))
}

func Csv2Multiple(inlet string) {
    C.csv2multiple(C.CString(inlet))
}
