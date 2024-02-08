package edf

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

/* --- MAIN FUNCTIONS --- */

// WriteGo writes the read data as Go vars.
func (edf *Edf) WriteGo() {
	fmt.Printf("header: %#v\n\n", edf.Header)
	fmt.Printf("records: %#v\n", edf.Records)
}

// WriteCSV formats the data to the *.csv format into a string.
// Ignores the annotations channel.
func (edf *Edf) WriteCSV() string {
	var buffer bytes.Buffer
	numberSignals := getNumberSignals(edf.Header)
	notesChannel := getAnnotationsChannel(edf.Header)

	// writing header...
	recorded := fmt.Sprintf("recorded:%s %s;",
		edf.Header["startdate"],
		edf.Header["starttime"])
	sampling := fmt.Sprintf("sampling:%v;", edf.GetSampling())
	patient := fmt.Sprintf("subject:%s;", edf.Header["patient"])
	labels := fmt.Sprintf("labels:%s;", strings.Join(edf.GetLabels(), ""))
	channel := fmt.Sprintf("chan:%s;", edf.Header["numbersignals"])
	units := fmt.Sprintf("units:%s\n", edf.GetUnits())

	buffer.WriteString(recorded)
	buffer.WriteString(sampling)
	buffer.WriteString(patient)
	buffer.WriteString(labels)
	buffer.WriteString(channel)
	buffer.WriteString(units)

	// writing data edf.Records...
	limit := len(edf.Records[0])
	for j := 0; j < limit; j++ {
		line := ""
		for i := 0; i < numberSignals; i++ {
			if i != notesChannel {
				data := edf.PhysicalRecords[i][j]

				if i == 0 {
					line += fmt.Sprintf("%f", data)
				} else {
					line += fmt.Sprintf("; %f", data)
				}
			}
		}
		buffer.WriteString(line + "\n")
	}

	outlet := buffer.String()
	return outlet
}

// WriteCsvToFile writes a CSV string directly to a file
func (edf *Edf) WriteCsvToFile(output string) {
	fp, _ := os.Create(output)
	buffer := bufio.NewWriter(fp)
	numberSignals := getNumberSignals(edf.Header)
	notesChannel := getAnnotationsChannel(edf.Header)
	defer fp.Close()

	// writing header...
	recorded := fmt.Sprintf("recorded:%s %s;",
		edf.Header["startdate"],
		edf.Header["starttime"])
	sampling := fmt.Sprintf("sampling:%v;", edf.GetSampling())
	patient := fmt.Sprintf("subject:%s;", edf.Header["patient"])
	labels := fmt.Sprintf("labels:%s;", strings.Join(edf.GetLabels(), ""))
	channel := fmt.Sprintf("chan:%s;", edf.Header["numbersignals"])
	units := fmt.Sprintf("units:%s\n", edf.GetUnits())

	buffer.WriteString(recorded)
	buffer.WriteString(sampling)
	buffer.WriteString(patient)
	buffer.WriteString(labels)
	buffer.WriteString(channel)
	buffer.WriteString(units)
	buffer.Flush()

	// writing data edf.Records...
	limit := len(edf.Records[0])
	for j := 0; j < limit; j++ {
		line := ""
		for i := 0; i < numberSignals; i++ {
			if i != notesChannel {
				data := edf.PhysicalRecords[i][j]

				if i == 0 {
					line += fmt.Sprintf("%f", data)
				} else {
					line += fmt.Sprintf("; %f", data)
				}
			}
		}
		buffer.WriteString(line + "\n")
		buffer.Flush()
	}

}

// WriteASCII translates the data to the *.ascii format into a string.
// Ignores the annotations channel.
func (edf *Edf) WriteASCII() string {
	numberSignals := getNumberSignals(edf.Header)
	notesChannel := getAnnotationsChannel(edf.Header)
	outlet := ""
	flag := numberSignals
	j := 0 // line number

	for flag > 0 {
		flag = 0

		for i := 0; i < numberSignals; i++ {
			if i != notesChannel {
				data, count := writeASCIIChannel(edf.PhysicalRecords[i], j)
				outlet += data
				flag += count
			}
		}

		outlet += fmt.Sprintf("\n")
		j++
	}

	return outlet
}

// WriteNotes extracts the annoatations channel from the EDF file, if it exists.
func (edf *Edf) WriteNotes() string {
	which := getAnnotationsChannel(edf.Header)
	outlet := ""

	if which >= 0 && which < len(edf.Records) {
		annotations := convertInt16ToByte(edf.Records[which])
		outlet += fmt.Sprintf("%s\n", formatAnnotations(annotations))
	}

	return outlet
}

// Get the labels' names from the EDF file in one String
func getLabels(header map[string]string) string {
	numberSignals := getNumberSignals(header)
	labels := separateString(header["label"], numberSignals)
	outlet := ""

	for i := 1; i < numberSignals; i++ {
		outlet += labels[i] + " "
	}

	return outlet
}

func getAnnotationsChannel(header map[string]string) int {
	result := -1
	labels := separateString(header["label"], getNumberSignals(header))

	for i, label := range labels {
		if match(label, "EDF Annotations") {
			result = i
		}
	}

	return result
}

/* returns false when it can't write anymore */
func writeASCIIChannel(record []float64, index int) (string, int) {
	outlet := ""
	flag := 1

	if index < len(record) {
		outlet += fmt.Sprintf("%f ", record[index])
	} else {
		outlet += fmt.Sprintf("0 ")
		flag = 0
	}

	return outlet, flag
}

/* format annotations to human-readable text */
func formatAnnotations(raw []byte) string {
	return faf(0, raw, false, "")
}

// aka formatAnnotationsFeedback
func faf(index int, raw []byte, inside bool, box string) string {
	if index == len(raw) {
		return box
	} else if inside {
		if raw[index] == 0 {
			return faf(index+1,
				raw,
				false,
				box+fmt.Sprintf("\n"))
		}
		if raw[index] == 20 || raw[index] == 21 {
			raw[index] = ' '
		}
		return faf(index+1,
			raw,
			inside,
			box+fmt.Sprintf("%c", raw[index]))
	} else if raw[index] == '+' || raw[index] == '-' {
		return faf(index+1,
			raw,
			true,
			box+string(raw[index]))
	} else {
		return faf(index+1, raw, inside, box)
	}
}

/************
 * EDF SAVE *
 ************/

// WriteEdf writes the EDF data to the file whose name is the output string. This
// function is in experimental state and must be used carefully!
// BUG: Adds one second of empty data before and after the recording
func (edf *Edf) WriteEdf(output string) {
	fp, oops := os.Create(output)

	if oops != nil {
		panic(oops)
	} else {
		defer fp.Close()
	}

	// Writting header
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
			field = EnforceSize(field, specsLength[spec])
			fmt.Fprintf(fp, "%s", field)
		}
	}

	numberSignals := getNumberSignals(edf.Header)
	for j := index; j < limit; j++ {
		spec := specsList[j]
		field := edf.Header[spec]
		field = EnforceSize(field, specsLength[spec]*numberSignals)
		fmt.Fprintf(fp, "%s", field)
	}

	// Writting data records
	dataRecords := str2int(edf.Header["datarecords"])
	numberSamples := getNumberSamples(edf.Header)

	buffer := new(bytes.Buffer)
	for d := 0; d < dataRecords; d++ {
		for i := 0; i < numberSignals; i++ {
			lowerLimit := d * numberSamples[i]
			upperLimit := (d + 1) * numberSamples[i]
			record := edf.Records[i][lowerLimit:upperLimit]
			for _, value := range record {
				binary.Write(buffer, binary.LittleEndian, value)
			}
		}
	}
	fp.Write(buffer.Bytes())
}
