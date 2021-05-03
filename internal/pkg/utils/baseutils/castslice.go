package baseutils

import "strconv"

func StringsToUint(lines []string) ([]uint, error) {
	res := make([]uint, len(lines))
	for i, line := range lines {
		val, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = uint(val)
	}
	return res, nil
}
