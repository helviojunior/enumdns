package islazy

import (
	"strconv"
)

// LeftTrucate a string if its more than max
func LeftTrucate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[max:]
}

func FormatInt(n int) string {
	return FormatInt64(int64(n))
}

func FormatInt64(n int64) string {
    in := strconv.FormatInt(n, 10)
    numOfDigits := len(in)
    if n < 0 {
        numOfDigits-- // First character is the - sign (not a digit)
    }
    numOfCommas := (numOfDigits - 1) / 3

    out := make([]byte, len(in)+numOfCommas)
    if n < 0 {
        in, out[0] = in[1:], '-'
    }

    for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
        out[j] = in[i]
        if i == 0 {
            return string(out)
        }
        if k++; k == 3 {
            j, k = j-1, 0
            out[j] = '.'
        }
    }
}