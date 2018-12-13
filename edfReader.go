package edf

import (
	"bytes"
	"encoding/binary"
	"os"
)

/* --- MAIN FUNCTIONS --- */

// ReadFile reads an EDF file, parsing it into the header and the records.
// The header will be a map relating the properties to a string with the raw data
// The records will be a matrix storing the raw bytes in the file
func ReadFile(input string) Edf {
	inlet, err := os.Open(input)

	if err != nil {
		panic(err)
	}

	defer inlet.Close()
	header := ReadHeader(inlet)
	records := ReadRecords(inlet, header)

	return NewEdf(header, records)
}

// ReadHeader reads the header of an EDF file. Requires the opened EDF file, the list of
// specifications and the length in bytes for each of them, as described by
// the EDF standard. The specs can be accessed though the GetSpecsList
// function, and their lenghts through the GetSpecsLength one.
func ReadHeader(inlet *os.File) map[string]string {
	specsList := GetSpecsList()
	specsLength := GetSpecsLength()
	header := make(map[string]string)
	index := 0

	// Reading header's header
	for index < len(specsList) {
		spec := specsList[index]

		if spec == "label" {
			break
		} else {
			data := make([]byte, specsLength[spec])
			n, _ := inlet.Read(data)
			header[spec] = string(data[:n])
		}

		index++
	}

	// Reading header's records
	numberSignals := getNumberSignals(header)
	for j := index; j < len(specsList); j++ {
		spec := specsList[j]
		data := make([]byte, specsLength[spec]*numberSignals)
		n, _ := inlet.Read(data)
		header[spec] = string(data[:n])
	}

	return header
}

// ReadRecords reads the data records from the EDF file. Its parameters are the pointer to
// file; and header information, as returned by the ReadHeader function.
func ReadRecords(inlet *os.File, header map[string]string) [][]int16 {
	numberSignals := getNumberSignals(header)
	numberSamples := getNumberSamples(header)
	records := make([][]int16, numberSignals)
	dataRecords := str2int(header["datarecords"])

	// setup records
	for i := 0; i < numberSignals; i++ {
		records[i] = make([]int16, 0, dataRecords*numberSamples[i])
	}

	// Reading records
	dataRecordsSize := 2 * dataRecords * Sigma(numberSamples)
	data := make([]byte, dataRecordsSize)
	inlet.Read(data)
	transData := translate(data)

	// translate data
	i := 0
	for d := 0; d < dataRecords; d++ {
		for s := 0; s < numberSignals; s++ {
			step := numberSamples[s]
			piece := transData[i : i+step]
			records[s] = append(records[s], piece...)
			i += step
		}
	}

	return records
}

/* --- AUXILIAR FUNCTIONS --- */
func translate(inlet []byte) []int16 {
	var data int16
	limit := len(inlet) / 2
	outlet := make([]int16, limit)
	buffer := bytes.NewReader(inlet)

	for i := 0; i < limit; i++ {
		oops := binary.Read(buffer, binary.LittleEndian, &data)
		if oops == nil {
			outlet[i] = data
		} else {
			panic(oops)
		}
	}

	return outlet
}

func getNumberSignals(header map[string]string) int {
	raw := header["numbersignals"]
	return str2int(raw)
}

func getNumberSamples(header map[string]string) []int {
	numberSignals := getNumberSignals(header)
	numberSamples := make([]int, numberSignals)
	samples := separateString(header["samplesrecord"], numberSignals)

	for i := 0; i < numberSignals; i++ {
		numberSamples[i] = str2int(samples[i])
	}

	return numberSamples
}
