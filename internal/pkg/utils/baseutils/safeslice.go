package baseutils

func SafePage(arrLen, offset, limit int) (left, right int) {
	left = Min(offset, arrLen)
	right = Min(offset+limit, arrLen)
	return
}
