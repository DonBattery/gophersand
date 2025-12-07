// draw.go provides functions for drawing simple shapes on an image
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	DefaultFont *text.GoTextFaceSource
)

// Px sets the color of a single pixel
func Px(target *ebiten.Image, x, y int, c color.Color) {
	target.Set(x, y, c)
}

// Line draws a line between two points
func Line(target *ebiten.Image, x1, y1, x2, y2 int, c color.Color) {
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
		target.Set(x1, y1, c)

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
}

// Rect draws a rectangle
func Rect(target *ebiten.Image, x, y, w, h int, c color.Color) {
	if w < 1 || h < 1 {
		return
	}

	for i := y; i < y+h; i++ {
		target.Set(x, i, c)
		target.Set(x+w-1, i, c)
	}

	for i := x; i < x+w; i++ {
		target.Set(i, y, c)
		target.Set(i, y+h-1, c)
	}
}

// Tx draws text
func Tx(target *ebiten.Image, msg string, x, y int, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(c)
	text.Draw(target, msg, &text.GoTextFace{
		Source: DefaultFont,
		Size:   5,
	}, op)
}
