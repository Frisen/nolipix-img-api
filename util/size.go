package util

//计算等比例缩放大小
func CalXY(srcX, srcY, dstX, dstY float64) (x int, y int) {
	if dstX == 0 && dstY == 0 {
		return int(srcX), int(srcY)
	}
	if dstX == 0 {
		x = int(srcX * (dstY / srcY))
		y = int(srcY * (dstY / srcY))
		return
	}
	if dstY == 0 {
		x = int(srcX * (dstX / srcX))
		y = int(srcY * (dstX / srcX))
		return
	}
	if dstX/srcX <= dstY/srcY {
		x = int(srcX * (dstX / srcX))
		y = int(srcY * (dstX / srcX))
		return
	}
	if dstX/srcX > dstY/srcY {
		x = int(srcX * (dstY / srcY))
		y = int(srcY * (dstY / srcY))
		return
	}
	return
}
