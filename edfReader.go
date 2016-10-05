package edf

import "os"
import "bytes"
import "encoding/binary"

/* --- MAIN FUNCTIONS --- */

// Reads and EDF file, parsing it into the header and the records.
// The header will be a map relating the properties to a string with the raw data
// The records will be a matrix storing the raw bytes in the file
func ReadFile(input string) (map[string]string, [][]int16) {
    inlet, _ := os.Open(input)
    specsList := GetSpecsList()
    specsLength := GetSpecsLength()

    defer inlet.Close()
    header := ReadHeader(inlet, specsList, specsLength)
    records := ReadRecords(inlet, header)

    return header, records
}

// Reads the header of an EDF file. Requires the opened EDF file, the list of
// specifications and the length in bytes for each of them, as described by
// the EDF standard. The specs can be accessed though the GetSpecsList
// function, and their lenghts through the GetSpecsLength one.
func ReadHeader(inlet *os.File,
                specsList []string,
                specsLength map[string]int) map[string]string {
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
    for index = index; index < len(specsList); index++ {
        spec := specsList[index]
        data := make([]byte, specsLength[spec] * numberSignals)
        n, _ := inlet.Read(data)
        header[spec] = string(data[:n])
    }

    return header
}

// Reads the data records from the EDF file. Its parameters are the pointer to
// file; and header information, as returned by the ReadHeader function.
func ReadRecords(inlet *os.File, header map[string]string) [][]int16 {
    numberSignals := getNumberSignals(header)
    numberSamples := getNumberSamples(header)
    records := make([][]int16, numberSignals)
    sampling := make([]int, numberSignals)
    duration := str2int(header["duration"])
    dataRecords := str2int(header["datarecords"])

    // setup records
    for i := 0; i < numberSignals; i++ {
        sampling[i] = duration * numberSamples[i]
        records[i] = make([]int16, sampling[i])
    }

    // translate data
    for d := 0; d < dataRecords; d++ {
        for i := 0; i < numberSignals; i++ {
            data := make([]byte, 2*sampling[i])
            inlet.Read(data)
            records[i] = append(records[i], translate(data))
        }
    }

    return records
}


/* --- AUXILIAR FUNCTIONS --- */
func translate(inlet []byte) []int16 {
    var data int16
    limit := len(inlet)/2
    outlet := make([]int16, limit)
    buffer := bytes.NewReader(inlet)

    for i := 0; i < limit; i++ {
        // oops := binary.Read(buffer, binary.BigEndian, &data)
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
    return str2int(header["numbersignals"])
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
