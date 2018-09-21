package edf

import (
	"errors"
)

// /!\ EXPERIMENTAL FUNCTION /!\
// 
// Appends data from one EDF to another. Returns a new EDF object and an error,
// which is `nil` if everything runs ok. This function requires the 
// EDF files to have:
//
// - Same sampling rate
// - Same number of channels
// - All units as equal
//
func Append(x, y Edf) (*Edf, error) {
	z := NewEdf(x.header, x.records)
	oops := nil
	
	// TODO Append y records to z
	// TODO Append update header info to be consistent with updates
	
	return z, oops
}