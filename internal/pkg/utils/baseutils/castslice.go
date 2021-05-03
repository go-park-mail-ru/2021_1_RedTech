package baseutils

import "strconv"

func StringsToUint(lines []string) ([]uint, error) {
	res := make([]uint, len(lines))
	for i, value := range lines {
		res[i] := strconv.ParseUint(value, 10, 64)
	}
	return res, nil
}
