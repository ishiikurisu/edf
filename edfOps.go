package edf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/******************
 * EDF PROPERTIES *
 ******************/

// GetDuration gets the duration of the file in seconds
func (edf Edf) GetDuration() float64 {
	v, oops := strconv.ParseFloat(strings.TrimSpace(edf.Header["duration"]), 64)
	if oops != nil {
		panic(oops)
	}
	return v
}

// GetNumberSignals gets the number of signals that are present in the EDF file
func (edf Edf) GetNumberSignals() int {
	return getNumberSignals(edf.Header)
}

// GetNumberSamples gets the number of samples in each channel
func (edf Edf) GetNumberSamples() []int {
	return getNumberSamples(edf.Header)
}

// GetLabels gets the labels' names from the EDF file in one array
func (edf Edf) GetLabels() []string {
	rawLabels := separateString(edf.Header["label"], getNumberSignals(edf.Header))
	limit := len(rawLabels)
	labels := make([]string, limit)

	for i, rawLabel := range rawLabels {
		labels[i] = strings.Replace(strings.Replace(rawLabel, "\n", " ", -1), "\r", " ", -1)
	}

	return labels
}

// GetConvertionFactors gets the convertion factor to each channel.
func (edf Edf) GetConvertionFactors() []float64 {
	ns := getNumberSignals(edf.Header)
	factors := make([]float64, ns)
	dmaxs := separateString(edf.Header["digitalmaximum"], ns)
	dmins := separateString(edf.Header["digitalminimum"], ns)
	pmaxs := separateString(edf.Header["physicalmaximum"], ns)
	pmins := separateString(edf.Header["physicalminimum"], ns)

	for i := 0; i < ns; i++ {
		dmax := str2float64(dmaxs[i])
		dmin := str2float64(dmins[i])
		pmax := str2float64(pmaxs[i])
		pmin := str2float64(pmins[i])
		dig := dmax - dmin
		phi := pmax - pmin
		factors[i] = dig / phi
	}

	return factors
}

// GetUnits gets the physical units from the recording.
// TODO extract units
func (edf Edf) GetUnits() string {
	return "uV"
}

// GetSampling gets the sampling rate from the recording.
func (edf Edf) GetSampling() int {
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

	return rates[0]
}

// GetDataRecords gets the number of data records
func (edf Edf) GetDataRecords() int {
	return str2int(edf.Header["datarecords"])
}

/***************
 * EDF METHODS *
 ***************/

// Append appends data from one EDF to another. Returns a new EDF object and an error,
// which is `nil` if everything runs ok. This function requires the
// EDF files to have the same sampling rate, the same number of channels,
// and the same duration for each data record.
//
// This function is in experimental state and must be used carefully!
func Append(x, y Edf) (*Edf, error) {
	// Checking for viability
	if x.GetNumberSignals() != y.GetNumberSignals() {
		return nil, errors.New("EDF files don't have the same number of records")
	}
	if x.GetSampling() != y.GetSampling() {
		return nil, errors.New("EDF files don't have the same sampling rate")
	}
	if x.GetDuration() != y.GetDuration() {
		return nil, errors.New("EDF files don't have the same sampling duration")
	}
	z := NewEdf(x.Header, x.Records, GetConvertedRecords(&(x.Records), x.Header))

	// Updating header
	z.Header["datarecords"] = EnforceSize(strconv.Itoa(x.GetDataRecords()+y.GetDataRecords()), GetSpecsLength()["datarecords"])

	// Appending data records
	for i := 0; i < len(x.Records); i++ {
		z.Records[i] = appendInt16Arrays(x.Records[i], y.Records[i])
	}

	return &z, nil
}
