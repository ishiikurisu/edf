package edf

import "math"

func elvis(pred bool, a int16, b int16) int16 {
	if pred {
		return a
	} else {
		return b
	}
}

func identifyOverflow(inlet []int16) []int16 {
	var last int16 = 0
    limit := len(inlet)
    outlet := make([]int16, limit)

    for i := 0; i < limit; i++ {
        it := inlet[i]
        diff := float64(it) - float64(last)

        if math.Abs(diff) > 16000 { // half int16 is a hell of a diff
        	outlet[i] = elvis(diff > 0, 500, -500)
        } else {
        	outlet[i] = 0
        }

        last = it
    }

    return outlet
}

func rearrange(inlet []int16) []int16 {
    limit := len(inlet)
    overflow := identifyOverflow(inlet)
    midlet := make([]int, limit)
    step := int(math.Pow(2, 15)) - 1
    factor := 0

    for i := 0; i < limit; i++ {
    	if overflow[i] < 0 {
    		factor += step
    	} else if overflow[i] > 0 {
    		factor -= step
    	}

        midlet[i] = int(inlet[i]) + factor
    }

    return convert(midlet)
}

func convert(midlet []int) []int16 {
	limit := len(midlet)
	outlet := make([]int16, limit)	

	for i := 0; i < limit; i++ {
		outlet[i] = int16(midlet[i]/2.0)
	}

	return outlet
}