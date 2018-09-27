package edf

import "fmt"
import "bytes"
import "strings"
import "bufio"
import "os"

/* --- MAIN FUNCTIONS --- */

// Just writes the read data as Go vars.
func (edf *Edf) WriteGo() {
	fmt.Printf("header: %#v\n\n", edf.Header)
	fmt.Printf("records: %#v\n", edf.Records)
}

// Formats the data to the *.csv format into a string.
// Ignores the annotations channel.
func (edf *Edf) WriteCSV() string {
	var buffer bytes.Buffer
	numberSignals := getNumberSignals(edf.Header)
	convertionFactor := edf.GetConvertionFactors()
	notesChannel := getAnnotationsChannel(edf.Header)

	// writing header...
	fmt.Sprintf("title:%s;", edf.Header["recording"])
	recorded := fmt.Sprintf("recorded:%s %s;",
		                     edf.Header["startdate"],
		                     edf.Header["starttime"])
	sampling := fmt.Sprintf("sampling:%s;", edf.GetSampling())
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
				data := float64(edf.Records[i][j]) * convertionFactor[i]

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

// Writes a CSV string directly to a file
func (edf *Edf) WriteCsvToFile(output string) {
	fp, _ := os.Create(output)
	buffer := bufio.NewWriter(fp)
	numberSignals := getNumberSignals(edf.Header)
	convertionFactor := edf.GetConvertionFactors()
	notesChannel := getAnnotationsChannel(edf.Header)
	defer fp.Close()

	// writing header...
	fmt.Sprintf("title:%s;", edf.Header["recording"])
	recorded := fmt.Sprintf("recorded:%s %s;",
		                     edf.Header["startdate"],
		                     edf.Header["starttime"])
	sampling := fmt.Sprintf("sampling:%s;", edf.GetSampling())
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
				data := float64(edf.Records[i][j]) * convertionFactor[i]

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

// Translates the data to the *.ascii format into a string.
// Ignores the annotations channel.
func (edf *Edf) WriteASCII() string {
	numberSignals := getNumberSignals(edf.Header)
	convertionFactor := edf.GetConvertionFactors()
	notesChannel := getAnnotationsChannel(edf.Header)
	outlet := ""
	flag := numberSignals
	j := 0 // line number

	for flag > 0 {
		flag = 0

		for i := 0; i < numberSignals; i++ {
			if i != notesChannel {
				data, count := writeASCIIChannel(edf.Records[i],
					                             convertionFactor[i],
					                             j)
				outlet += data
				flag += count
			}
		}

		outlet += fmt.Sprintf("\n")
		j += 1
	}

	return outlet
}

// Extracts the annoatations channel from the EDF file, if it exists.
func (edf *Edf) WriteNotes() string {
	which := getAnnotationsChannel(edf.Header)
	outlet := ""

	if which > 0 && which < len(edf.Records) {
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
func writeASCIIChannel(record []int16, factor float64, index int) (string, int) {
	outlet := ""
	flag := 1

	if index < len(record) {
		outlet += fmt.Sprintf("%f ", float64(record[index]) * factor)
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
			return faf(index + 1,
				       raw,
				       false,
				       box + fmt.Sprintf("\n"))
		} else {
			if raw[index] == 20 || raw[index] == 21 { raw[index] = ' ' }
			return faf(index + 1,
				       raw,
				       inside,
				       box + fmt.Sprintf("%c", raw[index]))
		}
	} else if raw[index] == '+' || raw[index] == '-' {
		return faf(index + 1,
				   raw,
				   true,
				   box + string(raw[index]))
	} else {
		return faf(index+1, raw, inside, box)
	}
}
