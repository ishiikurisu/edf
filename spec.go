package edf

/**************
 * EDF STRUCT *
 **************/

// Edf is the definition of the EDF structure to be used by this library.
type Edf struct {
	// The variable to hold the EDF's header information.
	Header map[string]string

	// The records will be stored in its raw form, each one of them stored in
	// an array of int16.
	Records [][]int16

	// additionally all records will be stored in its physical form
	PhysicalRecords [][]float64
}

// NewEdf returns a new Edf struct
func NewEdf(header map[string]string, records [][]int16, physicalRecords [][]float64) Edf {
	return Edf{
		Header:          header,
		Records:         records,
		PhysicalRecords: physicalRecords,
	}
}

/****************
 * SPECS LENGTH *
 ****************/

// GetSpecsLength gets the length in bytes of every specified field in the EDF file's header.
func GetSpecsLength() map[string]int {
	spec := make(map[string]int)

	spec["version"] = 8
	spec["patient"] = 80
	spec["recording"] = 80
	spec["startdate"] = 8
	spec["starttime"] = 8
	spec["bytesheader"] = 8
	spec["reserved"] = 44
	spec["datarecords"] = 8
	spec["duration"] = 8
	spec["numbersignals"] = 4
	spec["label"] = 16
	spec["transducer"] = 80
	spec["physicaldimension"] = 8
	spec["physicalminimum"] = 8
	spec["physicalmaximum"] = 8
	spec["digitalminimum"] = 8
	spec["digitalmaximum"] = 8
	spec["prefiltering"] = 80
	spec["samplesrecord"] = 8
	spec["chanreserved"] = 32

	return spec
}

// GetSpecsList gets the a list with codes for every field specified in the EDF file's
// header. They will appear in the order they are needed.
func GetSpecsList() []string {
	spec := make([]string, 20)

	spec[0] = "version"
	spec[1] = "patient"
	spec[2] = "recording"
	spec[3] = "startdate"
	spec[4] = "starttime"
	spec[5] = "bytesheader"
	spec[6] = "reserved"
	spec[7] = "datarecords"
	spec[8] = "duration"
	spec[9] = "numbersignals"
	spec[10] = "label"
	spec[11] = "transducer"
	spec[12] = "physicaldimension"
	spec[13] = "physicalminimum"
	spec[14] = "physicalmaximum"
	spec[15] = "digitalminimum"
	spec[16] = "digitalmaximum"
	spec[17] = "prefiltering"
	spec[18] = "samplesrecord"
	spec[19] = "chanreserved"

	return spec
}
