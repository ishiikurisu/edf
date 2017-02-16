package edf

import "os"
import "fmt"

/* ##################
   # MAIN FUNCTIONS #
   ################## */

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
	specsLength := GetSpecsLength()
	limit := len(specsList)

	for index := 0; index < limit; index++ {
		spec := specsList[index]

		if spec == "label" {
			break
		} else {
			// TODO Enforce length of header data
			field := edf.Header[spec]
			field = enforceSize(field, specsLength[spec])
			fmt.Fprintf(fp, "%s", field)
			fmt.Printf("%s: %s\n", spec, field)
		}
	}
	fmt.Println("")
}


/* ######################
   # AUXILIAR FUNCTIONS #
   ###################### */

func enforceSize(field string, limit int) string {
	length := len(field)

	if limit == length {

	} else if length > limit {
		field = field[0:limit]
	} else {
		for i := 0; i < limit - length; i++ {
			field += " "
		}
	}

	return field
}

/* CODE IS POETRY */
