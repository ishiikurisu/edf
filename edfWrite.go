package edf

import "fmt"
import "bytes"
import "strings"

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

/* --- AUXILIAR FUNCTIONS --- */

// Gets the sampling rate from the recording.
func (edf *Edf) GetSampling() string {
	ns := getNumberSignals(edf.Header)
	raw := separateString(edf.Header["samplesrecord"], ns)
	rates := make([]int, ns)

	// Turning sampling rates into numbers
	for i := 0; i < ns; i++ {
		fmt.Sscanf(raw[i], "%d", &rates[i])
	}

	// Getting most common designated sampling rate
	// TODO Write this part too
	// After thought: this might not be needed

	outlet := fmt.Sprintf("%d", rates[0])
	return outlet
}

// Gets the physical units from the recording.
func (edf *Edf) GetUnits() string {
	// TODO extract units
	return "uV"
}

// Gets the convertion factor to each channel.
func (edf *Edf) GetConvertionFactors() []float64 {
	ns := getNumberSignals(edf.Header)
	factors := make([]float64, ns)
	dmaxs := separateString(edf.Header["digitalmaximum"], ns)
	dmins := separateString(edf.Header["digitalminimum"], ns)
	pmaxs := separateString(edf.Header["physicalmaximum"], ns)
	pmins := separateString(edf.Header["physicalminimum"], ns)

	for i := 0; i < ns; i++ {
		dmax := str2int64(dmaxs[i])
		dmin := str2int64(dmins[i])
		pmax := str2int64(pmaxs[i])
		pmin := str2int64(pmins[i])
		dig := float64(dmax-dmin)
		phi := float64(pmax-pmin)
		factors[i] = dig/phi;
	}

	return factors
}

// Get the labels' names from the EDF file in one array
func (edf *Edf) GetLabels() []string {
	return separateString(edf.Header["label"], getNumberSignals(edf.Header))
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
	return formatAnnotationsFeedback(0, raw, false, "")
}
func formatAnnotationsFeedback(index int, raw []byte, inside bool, box string) string {
	if index == len(raw) {
		return box
	} else if inside {
		if raw[index] == 0 {
			return formatAnnotationsFeedback(index + 1,
				                             raw,
				                             false,
				                             box + fmt.Sprintf("\n"))
		} else {
			if raw[index] == 20 || raw[index] == 21 { raw[index] = ' ' }
			return formatAnnotationsFeedback(index + 1,
				                             raw,
				                             inside,
				                             box + fmt.Sprintf("%c", raw[index]))
		}
	} else if raw[index] == '+' || raw[index] == '-' {
		return formatAnnotationsFeedback(index + 1,
				                         raw,
				                         true,
				                         box + string(raw[index]))
	} else {
		return formatAnnotationsFeedback(index+1, raw, inside, box)
	}
}
