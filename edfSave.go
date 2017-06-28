package edf

import "os"
import "fmt"
import "bytes"
import "encoding/binary"

/* ##################
   # MAIN FUNCTIONS #
   ################## */

// Writes the EDF data to the file whose name is the output string. This
// function is in experimental state and must be used carefully!
// BUG: Adds one second of empty data before and after the recording
func (edf *Edf) WriteEdf(output string) {
	fp, oops := os.Create(output)

	if oops != nil {
		panic(oops)
	} else {
		defer fp.Close()
	}

	writeHeaderHeader(edf, fp)
	writeRecords(edf, fp)
}

func writeHeaderHeader(edf *Edf, fp *os.File) {
	specsList := GetSpecsList()
	specsLength := GetSpecsLength()
	limit := len(specsList)
	index := 0

	for index = 0; index < limit; index++ {
		spec := specsList[index]

		if spec == "label" {
			break
		} else {
			field := edf.Header[spec]
			field = enforceSize(field, specsLength[spec])
			fmt.Fprintf(fp, "%s", field)
			fmt.Printf("%s: %s\n", spec, field)
		}
	}

	numberSignals := getNumberSignals(edf.Header)
	for index = index; index < limit; index++ {
		spec := specsList[index]
		field := edf.Header[spec]
		field = enforceSize(field, specsLength[spec] * numberSignals)
		fmt.Fprintf(fp, "%s", field)
		fmt.Printf("%s: %s\n", spec, field)
    }
}

func writeRecords(edf *Edf, fp *os.File) {
	dataRecords := str2int(edf.Header["datarecords"])
	numberSignals := getNumberSignals(edf.Header)
	sampling := make([]int, numberSignals)
	duration := str2int(edf.Header["duration"])
	numberSamples := getNumberSamples(edf.Header)

	// Preparing records
	for i := 0; i < numberSignals; i++ {
		sampling[i] = duration * numberSamples[i]
	}

	// Writting chops
	for d := 0; d < dataRecords-1; d++ {
		for i := 0; i < numberSignals; i++ {
			lowerLimit := d * sampling[i]
			upperLimit := (d+1) * sampling[i]
			record := edf.Records[i][lowerLimit:upperLimit]
			for _, value := range record {
				buffer := new(bytes.Buffer)
				binary.Write(buffer, binary.LittleEndian, value)
				fp.Write(buffer.Bytes())
			}
		}
	}
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
