package edf

import "fmt"
import "math"
import "strconv"
import "strings"

func separateString(stuff string, howMany int) []string {
	bit := len(stuff) / howMany
	outlet := make([]string, howMany)

	for i := 0; i < howMany; i++ {
		outlet[i] = stuff[bit*i : bit*(i+1)]
	}

	return outlet
}

func appendInt16Arrays(original, toAppend []int16) []int16 {
	lo := len(original)
	lt := len(toAppend)
	outlet := make([]int16, lo+lt)

	for o := 0; o < lo; o++ {
		outlet[o] = original[o]
	}
	for t := 0; t < lt; t++ {
		outlet[lo+t] = toAppend[t]
	}

	return outlet
}
func str2int(inlet string) int {
	var outlet int
	fmt.Sscanf(inlet, "%d", &outlet)
	return outlet
}

func str2int64(str string) int64 {
	var x int64
	fmt.Sscanf(str, "%d", &x)
	return x
}
func str2float64(str string) float64 {
	x, oops := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if oops != nil {
		return float64(str2int64(str))
	}
	return x
}

func min(u, v int) int {
	if u < v {
		return u
	}
	return v
}

func match(p, q string) bool {
	r := true
	l := min(len(p), len(q))

	for i := 0; i < l && r; i++ {
		if p[i] != q[i] {
			r = false
		}
	}

	return r
}

func smallest(array []int) int {
	s := array[0]
	r := 0

	for i, it := range array {
		if it < s {
			s = it
			r = i
		}
	}

	return r
}

func convertInt16ToByte(inlet []int16) []byte {
	outlet := make([]byte, 2*len(inlet))
	var this byte
	var i uint
	j := 0

	for _, that := range inlet {
		bits := extractBits(that)

		this = 0
		for i = 0; i < 8; i++ {
			this += bits[i] * byte(math.Pow(2, float64(i)))
		}
		outlet[j] = this
		j++

		this = 0
		for i = 8; i < 16; i++ {
			this += bits[i] * byte(math.Pow(2, float64(i-8)))
		}
		outlet[j] = this
		j++
	}

	return outlet
}

func extractBits(inlet int16) []byte {
	outlet := make([]byte, 16)
	var i uint

	for i = 0; i < 16; i++ {
		outlet[i] = byte(inlet>>i) & 0x1
	}

	return outlet
}

// Iff is a functional IF statemant.
func Iff(s bool, t, e string) string {
	if s {
		return t
	}
	return e
}

// EnforceSize forces a string to be a certain size.
func EnforceSize(field string, limit int) string {
	length := len(field)

	if limit == length {

	} else if length > limit {
		field = field[0:limit]
	} else {
		for i := 0; i < limit-length; i++ {
			field += " "
		}
	}

	return field
}

// Sigma sums all values of an int array
func Sigma(a []int) int {
	r := 0
	for i := 0; i < len(a); i++ {
		r += a[i]
	}
	return r
}
