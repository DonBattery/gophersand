// utils.go provides general utility functions
package game

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// mid returns the middle value among a, b, and c
func mid(a, b, c int) int {
	if a > b {
		a, b = b, a
	}

	if b > c {
		b, c = c, b
	}

	if a > b {
		a, b = b, a
	}

	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func clamp(val, mi, ma int) int {
	return min(max(val, mi), ma)
}

func lerp(from, to int, t float64) int {
	return int(float64(from) + (float64(to-from) * t))
}

func minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func midf(a, b, c float32) float32 {
	if a > b {
		a, b = b, a
	}

	if b > c {
		b, c = c, b
	}

	if a > b {
		a, b = b, a
	}

	return b
}

func absf(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func clampf(val, mi, ma float32) float32 {
	return minf(maxf(val, mi), ma)
}

func ray(v1, v2 *Vector) []int {
	pts := []int{}

	x1, y1 := int(v1.X), int(v1.Y)
	x2, y2 := int(v2.X), int(v2.Y)
	dx := abs(x1 - x2)
	dy := abs(y1 - y2)

	sx := -1
	if x1 < x2 {
		sx = 1
	}

	sy := -1
	if y1 < y2 {
		sy = 1
	}

	err := dx - dy

	for {
		pts = append(pts, x1, y1)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := err * 2

		if e2 > -dy {
			err -= dy
			x1 += sx
		}

		if e2 < dx {
			err += dx
			y1 += sy
		}
	}

	return pts
}
