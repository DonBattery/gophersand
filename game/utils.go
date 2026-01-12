// utils.go provides general utility functions
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

func lerpf(from, to float32, t float64) float32 {
	return float32(float64(from) + (float64(to-from) * t))
}

func Line(target *ebiten.Image, x1, y1, x2, y2 int, c color.Color) {
	vector.StrokeLine(target, float32(x1), float32(y1), float32(x2), float32(y2), 1, c, false)
}

func Rect(target *ebiten.Image, x, y, w, h int, c color.Color) {
	if w < 1 || h < 1 {
		return
	}

	Line(target, x, y, x+w-1, y, c)
	Line(target, x, y+h-1, x+w-1, y+h-1, c)
	Line(target, x, y, x, y+h-1, c)
	Line(target, x+w-1, y, x+w-1, y+h-1, c)
}

func Circle(target *ebiten.Image, x, y, r int, c color.Color) {
	if r <= 1 {
		target.Set(x, y, c)
		return
	}
	vector.StrokeCircle(target, float32(x), float32(y), float32(r), 1, c, false)
}
