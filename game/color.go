// color.go provides functions for handling RGBA colors
package game

import (
	"fmt"
	"strconv"
)

// Color stores a pixel in Ebiten's native little-endian BGRA format.
type Color uint32

// ------------------------------------------------------------
// Construction
// ------------------------------------------------------------

// ColorRGBA builds a Color from standard RGBA components.
func ColorRGBA(r, g, b, a byte) Color {
	// Memory layout: [R, G, B, A] (RGBA little-endian)
	return Color(uint32(r) |
		uint32(g)<<8 |
		uint32(b)<<16 |
		uint32(a)<<24)
}

// ColorFromRGB has optional alpha.
func ColorFromRGB(r, g, b byte, alphas ...byte) Color {
	a := byte(255)
	if len(alphas) > 0 {
		a = alphas[0]
	}
	return ColorRGBA(r, g, b, a)
}

// ColorFromHex creates a color from "#RRGGBB".
func ColorFromHex(hex string, alpha ...byte) Color {
	r, _ := strconv.ParseUint(hex[1:3], 16, 8)
	g, _ := strconv.ParseUint(hex[3:5], 16, 8)
	b, _ := strconv.ParseUint(hex[5:7], 16, 8)

	a := byte(255)
	if len(alpha) > 0 {
		a = alpha[0]
	}

	return ColorRGBA(byte(r), byte(g), byte(b), a)
}

// WithAlpha replaces the alpha channel.
func (c Color) WithAlpha(a byte) Color {
	v := uint32(c)
	return Color((v & 0x00FFFFFF) | uint32(a)<<24)
}

// ------------------------------------------------------------
// Debug / Formatting
// ------------------------------------------------------------
func (c Color) String() string {
	r, g, b, a := c.UnpackRGBA8()
	return fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a)
}

// ------------------------------------------------------------
// Internal helpers
// ------------------------------------------------------------

// UnpackRGBA8 gets the original 8-bit RGBA values (non-premultiplied)
func (c Color) UnpackRGBA8() (r, g, b, a byte) {
	v := uint32(c)
	r = byte(v)       // lowest byte is now R
	g = byte(v >> 8)  // second byte is G
	b = byte(v >> 16) // third byte is B
	a = byte(v >> 24) // highest byte is A
	return
}

// ------------------------------------------------------------
// color.Color interface
// Ebiten expects RGBA64, premultiplied alpha.
// ------------------------------------------------------------
func (c Color) RGBA() (r, g, b, a uint32) {
	// Extract 8-bit channels (already in RGBA order)
	r8, g8, b8, a8 := c.UnpackRGBA8()

	// Convert to 16-bit
	r = uint32(r8) * 0x101
	g = uint32(g8) * 0x101
	b = uint32(b8) * 0x101
	a = uint32(a8) * 0x101

	// Premultiply (required by color.Color contract)
	if a > 0 {
		r = r * a / 0xFFFF
		g = g * a / 0xFFFF
		b = b * a / 0xFFFF
	}

	return
}

var (
	ColorNull  = ColorFromHex("#000000", 0)
	ColorBlack = ColorFromHex("#000000")
	ColorWhite = ColorFromHex("#ffffff")

	ColorActiveCell   = ColorFromHex("#58e075")
	ColorInactiveCell = ColorFromHex("#0e4c70ff")

	ColorRed        = ColorFromHex("#f53141")
	ColorLightRed   = ColorFromHex("#eb5f6f")
	ColorLighterRed = ColorFromHex("#eb8793")
	ColorDarkRed    = ColorFromHex("#9a0d27")
	ColorDarkerRed  = ColorFromHex("#61051a")

	ColorGreen        = ColorFromHex("#15c24e")
	ColorLightGreen   = ColorFromHex("#57db83")
	ColorLighterGreen = ColorFromHex("#8df397")
	ColorDarkGreen    = ColorFromHex("#0b8a4b")
	ColorDarkerGreen  = ColorFromHex("#074d3f")

	ColorBlue        = ColorFromHex("#5185df")
	ColorLightBlue   = ColorFromHex("#85b0eb")
	ColorLighterBlue = ColorFromHex("#84cdff")
	ColorDarkBlue    = ColorFromHex("#37559a")
	ColorDarkerBlue  = ColorFromHex("#1d2f55")

	ColorYellow        = ColorFromHex("#e69b22")
	ColorLightYellow   = ColorFromHex("#ffcd38")
	ColorLighterYellow = ColorFromHex("#f3e064")
	ColorDarkYellow    = ColorFromHex("#ce922a")
	ColorDarkerYellow  = ColorFromHex("#8a6b28")

	ColorPurple        = ColorFromHex("#a35dd9")
	ColorLightPurple   = ColorFromHex("#ca7ef2")
	ColorLighterPurple = ColorFromHex("#e29bfa")
	ColorDarkPurple    = ColorFromHex("#773bbf")
	ColorDarkerPurple  = ColorFromHex("#4e278c")

	ColorPink        = ColorFromHex("#e35c9b")
	ColorLightPink   = ColorFromHex("#f391ad")
	ColorLighterPink = ColorFromHex("#e7abc6")
	ColorDarkPink    = ColorFromHex("#b32d7d")
	ColorDarkerPink  = ColorFromHex("#852264")

	ColorGray        = ColorFromHex("#696570")
	ColorLightGray   = ColorFromHex("#807980")
	ColorLighterGray = ColorFromHex("#a69a9c")
	ColorDarkGray    = ColorFromHex("#495169")
	ColorDarkerGray  = ColorFromHex("#0d2140")

	ColorBrown        = ColorFromHex("#875d58")
	ColorLightBrown   = ColorFromHex("#9e7767")
	ColorLighterBrown = ColorFromHex("#b58c7f")
	ColorDarkBrown    = ColorFromHex("#6e4250")
	ColorDarkerBrown  = ColorFromHex("#472e3e")

	ColorLime        = ColorFromHex("#b8e325")
	ColorLightLime   = ColorFromHex("#ccef74")
	ColorLighterLime = ColorFromHex("#e2fba1")
	ColorDarkLime    = ColorFromHex("#91b239")
	ColorDarkerLime  = ColorFromHex("#506d19")

	ColorOrange        = ColorFromHex("#ef7a2c")
	ColorLightOrange   = ColorFromHex("#db8c56")
	ColorLighterOrange = ColorFromHex("#ebaa73")
	ColorDarkOrange    = ColorFromHex("#ba521b")
	ColorDarkerOrange  = ColorFromHex("#7d4230")

	EmptyColors = []Color{ColorNull}

	StoneColors = []Color{
		ColorFromHex("#1d2429"),
		ColorFromHex("#252e36"),
		ColorFromHex("#303b45"),
		ColorFromHex("#3d4b59"),
		ColorFromHex("#495b6b"),
	}

	SandColors = []Color{
		ColorFromHex("#e3b314"),
		ColorFromHex("#d1a515"),
		ColorFromHex("#ba9311"),
		ColorFromHex("#a38214"),
		ColorFromHex("#947716"),
	}

	WaterColors = []Color{
		ColorFromHex("#15a8d1"),
	}

	SmokeColors = []Color{
		ColorFromHex("#ada8b5"),
	}

	MaterialColorGroups = [][]Color{
		EmptyColors,
		StoneColors,
		SandColors,
		WaterColors,
		SmokeColors,
	}
)
