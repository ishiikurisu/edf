package edf

import (
	"errors"
	"fmt"
	"strconv"
)

// /!\ EXPERIMENTAL FUNCTION /!\
//
// Appends data from one EDF to another. Returns a new EDF object and an error,
// which is `nil` if everything runs ok. This function requires the
// EDF files to have:
//
// - Same sampling rate
// - Same number of channels
//
func Append(x, y Edf) (*Edf, error) {
	// TODO Check for viability
	if len(x.Records) != len(y.Records) {
		return nil, errors.New("EDF files don't have the same number of records")
	}
	z := NewEdf(x.Header, x.Records)

	// TODO Append y records to z
	for i := 0; i < len(x.Records); i++ {
		z.Records[i] = appendInt16Arrays(x.Records[i], y.Records[i])
	}

	z.Header["duration"] = enforceSize(strings.Itoa(x.GetDuration() + y.GetDuration()), 8)
	// TODO Update header field "number of samples"

	return &z, nil
}
