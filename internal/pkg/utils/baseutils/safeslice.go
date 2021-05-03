package baseutils

func SafePage(arrLen, offset, limit int) (left, right int) {
	left = Min(Max(0, offset), arrLen)
	right = Min(offset+Max(0, limit), arrLen)
	return
}

