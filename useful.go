package edf

import "fmt"
import "math"

func separateString(stuff string, howMany int) []string {
	bit := len(stuff) / howMany
    outlet := make([]string, howMany)

    for i := 0; i < howMany; i++ {
    	outlet[i] = stuff[bit*i:bit*(i+1)]
    }

    return outlet
}

func append(original, to_append []int16) []int16 {
    lo := len(original)
    lt := len(to_append)
    outlet := make([]int16, lo + lt)

    for o := 0; o < lo; o++ {
        outlet[o] = original[o]
    }
    for t := 0; t < lt; t++ {
        outlet[lo+t] = to_append[t]
    }

    return outlet
}
func str2int(inlet string) int {
    var outlet int = 0
    fmt.Sscanf(inlet, "%d", &outlet)
    return outlet
}

func str2int64(str string) int64 {
	var x int64
	fmt.Sscanf(str, "%d", &x)
	return x
}

func min(u, v int) int {
    if u < v {
        return u
    } else {
        return v
    }
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
    var outlet []byte = make([]byte, 2*len(inlet))
    var this byte = 0
    var i uint = 0
    var j int = 0

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
        outlet[i] = byte(inlet >> i) & 0x1
    }

    return outlet
}

// Deprecated. A functional IF statemant.
func Iff(s bool, t, e string) string {
    if s {
        return t
    } else {
        return e
    }
}
