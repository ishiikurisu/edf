package edf

import "os"
import "fmt"

// Writes the EDF data to the given provided by the output string.
// WARNING!! THIS FUNCTION IS EXPERIMENTAL.
func (edf *Edf) WriteEdf(output string) {
	fp, oops := os.Create(output)

	if oops != nil {
		panic(oops)
	} else {
		defer fp.Close()
	}

	// Saving header's header
	writeHeaderHeader(edf, fp)
	// TODO Save header's records
	// TODO Save records
}

func writeHeaderHeader(edf *Edf, fp *os.File) {
	specsList := GetSpecsList()
	limit := len(specsList)

	for index := 0; index < limit; index++ {
		spec := specsList[index]

		if spec == "label" {
			break
		} else {
			// TODO Enforce length of header data
			fmt.Fprintf(fp, "%s", edf.Header[spec])
			fmt.Printf("%s", edf.Header[spec])
		}
	}
	fmt.Println("")
}
