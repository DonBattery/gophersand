// color.go provides functions for handling RGBA colors
package game

import (
	"fmt"
	"strconv"
)

// Color represents an RGBA color
type Color uint32

// ColorRGBA builds a Color from standard RGBA components.
func ColorRGBA(r, g, b, a byte) Color {
	// Memory layout: [R, G, B, A] (RGBA little-endian)
	return Color(uint32(r) |
		uint32(g)<<8 |
		uint32(b)<<16 |
		uint32(a)<<24)
}

func ColorFromHex(hex string) Color {
	r, _ := strconv.ParseUint(hex[1:3], 16, 8)
	g, _ := strconv.ParseUint(hex[3:5], 16, 8)
	b, _ := strconv.ParseUint(hex[5:7], 16, 8)
	a := uint64(255)
	if len(hex) >= 9 { // "#RRGGBBAA" is length 9
		a, _ = strconv.ParseUint(hex[7:9], 16, 8)
	}
	return ColorRGBA(byte(r), byte(g), byte(b), byte(a))
}

// WithAlpha replaces the alpha channel.
func (c Color) WithAlpha(a byte) Color {
	v := uint32(c)
	return Color((v & 0x00FFFFFF) | uint32(a)<<24)
}

// String returns a string representation of the color.
func (c Color) String() string {
	r, g, b, a := c.UnpackRGBA8()
	return fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a)
}

// ------------------------------------------------------------
// Internal helpers
// ------------------------------------------------------------

// UnpackRGBA8 gets the original 8-bit RGBA values (non-premultiplied)
func (c Color) UnpackRGBA8() (r, g, b, a byte) {
	r = byte(c)       // lowest byte is R
	g = byte(c >> 8)  // second byte is G
	b = byte(c >> 16) // third byte is B
	a = byte(c >> 24) // highest byte is A
	return
}

// ------------------------------------------------------------
// color.Color interface
// RGBA implements color.Color.
// Ebiten expects premultiplied alpha.
// ------------------------------------------------------------
func (c Color) RGBA() (r, g, b, a uint32) {
	v := uint32(c)

	r8 := (v >> 0) & 0xFF
	g8 := (v >> 8) & 0xFF
	b8 := (v >> 16) & 0xFF
	a8 := (v >> 24) & 0xFF

	// Expand 8-bit -> 16-bit (0..65535)
	r = r8 * 0x101
	g = g8 * 0x101
	b = b8 * 0x101
	a = a8 * 0x101

	// Premultiply
	if a != 0 {
		r = r * a / 0xFFFF
		g = g * a / 0xFFFF
		b = b * a / 0xFFFF
	}

	return
}

var (
	ColorNull  = ColorFromHex("#00000000")
	ColorBlack = ColorFromHex("#000000")
	ColorWhite = ColorFromHex("#ffffff")

	ColorActiveCell   = ColorFromHex("#24eb4fb7")
	ColorInactiveCell = ColorFromHex("#3d5766b7")

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

	// ==========================================================================
	// Unified Material Color Array
	// ==========================================================================
	// Each material has 16 color entries:
	//   Index 0-3:   Normal colors (life 0-3)
	//   Index 4-7:   Burned colors (life 0-3)
	//   Index 8-11:  Acidic colors (life 0-3)
	//   Index 12-15: Frozen colors (life 0-3)
	//
	// Color index formula: (state * 4) + life
	// Where state: 0=Normal, 1=Burned, 2=Acidic, 3=Frozen
	//
	// Materials without variations (Water, Acid, Steam, Smoke) repeat the same
	// color for all entries in each state group.
	// ==========================================================================

	// MaterialColors is the flat color array, IndexedBy MaterialKind * 16 + Status * 4 + Life
	MaterialColors = [256]Color{
		// Empty (0) - transparent
		ColorNull, ColorNull, ColorNull, ColorNull, // Normal
		ColorNull, ColorNull, ColorNull, ColorNull, // Burned
		ColorNull, ColorNull, ColorNull, ColorNull, // Acidic
		ColorNull, ColorNull, ColorNull, ColorNull, // Frozen

		// Stone (1)
		ColorFromHex("#839bc6ff"), ColorFromHex("#6d7e9cff"), ColorFromHex("#4d5a6eff"), ColorFromHex("#2d3640ff"),
		ColorFromHex("#8f340cff"), ColorFromHex("#772a06ff"), ColorFromHex("#662606ff"), ColorFromHex("#300e03ff"),
		ColorFromHex("#279800ff"), ColorFromHex("#3a9100ff"), ColorFromHex("#508320ff"), ColorFromHex("#567560ff"),
		ColorFromHex("#7299e3ff"), ColorFromHex("#6183bfff"), ColorFromHex("#445c80ff"), ColorFromHex("#2c3c4fff"),

		// Sand (2)
		ColorFromHex("#9c8a15ff"), ColorFromHex("#c0b51aff"), ColorFromHex("#d4c64bff"), ColorFromHex("#ece760ff"),
		ColorFromHex("#21160fff"), ColorFromHex("#3d2518ff"), ColorFromHex("#57442fff"), ColorFromHex("#544e38ff"),
		ColorFromHex("#6c7312ff"), ColorFromHex("#98a117ff"), ColorFromHex("#9dc41eff"), ColorFromHex("#98ee27ff"),
		ColorFromHex("#8f782dff"), ColorFromHex("#be9c3eff"), ColorFromHex("#c5a13dff"), ColorFromHex("#f0c948ff"),

		// Water (3) - single color, no variations (New Golang blue)
		ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), // Normal
		ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), // Burned (same)
		ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), // Acidic (same)
		ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), ColorFromHex("#00ADD8"), // Frozen (same)

		// Seed (4)
		ColorFromHex("#414709ff"), ColorFromHex("#5d6b18ff"), ColorFromHex("#768e25ff"), ColorFromHex("#96c635ff"),
		ColorFromHex("#3e1f09ff"), ColorFromHex("#683c15ff"), ColorFromHex("#935922ff"), ColorFromHex("#e39055ff"),
		ColorFromHex("#174709ff"), ColorFromHex("#2a6b18ff"), ColorFromHex("#418e25ff"), ColorFromHex("#59c635ff"),
		ColorFromHex("#3f3b15ff"), ColorFromHex("#616126ff"), ColorFromHex("#908d40ff"), ColorFromHex("#a4ae56ff"),

		// Ant (5)
		ColorFromHex("#f3b8b2ff"), ColorFromHex("#973831ff"), ColorFromHex("#c83f30ff"), ColorFromHex("#f33737ff"),
		ColorFromHex("#973e0aff"), ColorFromHex("#813319ff"), ColorFromHex("#c65b2aff"), ColorFromHex("#de7e24ff"),
		ColorFromHex("#c3db4aff"), ColorFromHex("#a1b52dff"), ColorFromHex("#d1d72bff"), ColorFromHex("#edf57dff"),
		ColorFromHex("#d797deff"), ColorFromHex("#97544fff"), ColorFromHex("#cd594dff"), ColorFromHex("#f05858ff"),

		// Wasp (6)
		ColorFromHex("#f9ea3dff"), ColorFromHex("#cfa400ff"), ColorFromHex("#ffcf2eff"), ColorFromHex("#fff2a6ff"),
		ColorFromHex("#c07f3aff"), ColorFromHex("#3b1f00ff"), ColorFromHex("#5a2b00ff"), ColorFromHex("#7a3d00ff"),
		ColorFromHex("#cde245ff"), ColorFromHex("#ccf52eff"), ColorFromHex("#b3e622ff"), ColorFromHex("#8fd11aff"),
		ColorFromHex("#d6cf88ff"), ColorFromHex("#e3e9c5ff"), ColorFromHex("#eff3dbff"), ColorFromHex("#f7f9e6ff"),

		// Acid (7) - single color, no variations
		ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"),
		ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"),
		ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"),
		ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"), ColorFromHex("#1ff52aff"),

		// Fire (8) - uses life for intensity, uses state for variants
		ColorFromHex("#792911ff"), ColorFromHex("#c63e1cff"), ColorFromHex("#e38d54ff"), ColorFromHex("#ede19bff"), // Normal
		ColorFromHex("#763420ff"), ColorFromHex("#9b2e13ff"), ColorFromHex("#e77a31ff"), ColorFromHex("#e9d45bff"), // Burned
		ColorFromHex("#642613ff"), ColorFromHex("#b6310fff"), ColorFromHex("#cc6a29ff"), ColorFromHex("#ccb844ff"), // Acidic
		ColorFromHex("#471709ff"), ColorFromHex("#a0361bff"), ColorFromHex("#d26a25ff"), ColorFromHex("#dec011ff"), // Frozen

		// Ice (9) - uses life for melting state
		ColorFromHex("#225587ff"), ColorFromHex("#4f8dc7ff"), ColorFromHex("#7da4dcff"), ColorFromHex("#80c9e3ff"), // Normal
		ColorFromHex("#225587ff"), ColorFromHex("#4f8dc7ff"), ColorFromHex("#7da4dcff"), ColorFromHex("#80c9e3ff"), // Burned (same)
		ColorFromHex("#225587ff"), ColorFromHex("#4f8dc7ff"), ColorFromHex("#7da4dcff"), ColorFromHex("#80c9e3ff"), // Acidic (same)
		ColorFromHex("#225587ff"), ColorFromHex("#4f8dc7ff"), ColorFromHex("#7da4dcff"), ColorFromHex("#80c9e3ff"), // Frozen (same)

		// Smoke (10) - single color, no variations
		ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"),
		ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"),
		ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"),
		ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"), ColorFromHex("#817b70ff"),

		// Steam (11) - single color, no variations
		ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"),
		ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"),
		ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"),
		ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"), ColorFromHex("#88c8cfff"),

		// Root (12) - uses life for color intensity
		ColorFromHex("#c2923aff"), ColorFromHex("#ae8930ff"), ColorFromHex("#94712cff"), ColorFromHex("#5c480fff"), // Normal
		ColorFromHex("#7c5a24ff"), ColorFromHex("#644e2aff"), ColorFromHex("#544425ff"), ColorFromHex("#392e11ff"), // Burned (darker, more muted)
		ColorFromHex("#b8ac3aff"), ColorFromHex("#a09c3aff"), ColorFromHex("#888c3aff"), ColorFromHex("#67641cff"), // Acidic (more yellow/saturated)
		ColorFromHex("#a0b0c9ff"), ColorFromHex("#8c98abff"), ColorFromHex("#78868dff"), ColorFromHex("#5a635eff"), // Frozen (desaturated blue/gray tones)

		// Plant (13) - leaf/emerald greens (life 0..3: darkest -> brightest)
		ColorFromHex("#173012ff"), ColorFromHex("#23501bff"), ColorFromHex("#2f7424ff"), ColorFromHex("#49a83aff"), // Normal
		ColorFromHex("#231304ff"), ColorFromHex("#3d2308ff"), ColorFromHex("#5a360dff"), ColorFromHex("#7b5220ff"), // Burned (charred -> ashy brown)
		ColorFromHex("#2a3f10ff"), ColorFromHex("#3e5e14ff"), ColorFromHex("#5c861bff"), ColorFromHex("#86b92aff"), // Acidic (sickly yellow-green)
		ColorFromHex("#1f3430ff"), ColorFromHex("#2e4f49ff"), ColorFromHex("#3d6a63ff"), ColorFromHex("#5f8f88ff"), // Frozen (desaturated, slightly blue-green)

		// Flower (14) - Life field (0-3) = base color: Blue, Pink, Magenta, Yellow
		// Status field (0-3) = filter: Normal, Burned, Acidic, Frozen
		ColorFromHex("#1079e2ff"), ColorFromHex("#10e2bbff"), ColorFromHex("#e210d4ff"), ColorFromHex("#d7e210ff"), // Status 0 (Normal): Blue, Pink, Magenta, Yellow
		ColorFromHex("#7a4521ff"), ColorFromHex("#3d928aff"), ColorFromHex("#7a2d6fff"), ColorFromHex("#8a7421ff"), // Status 1 (Burned): darkened burnt versions
		ColorFromHex("#2ab8a0ff"), ColorFromHex("#29773fff"), ColorFromHex("#b828b8ff"), ColorFromHex("#b8b828ff"), // Status 2 (Acidic): slightly greenish/sickly tint
		ColorFromHex("#5a9dc7ff"), ColorFromHex("#43716eff"), ColorFromHex("#c75ac7ff"), ColorFromHex("#c7c75aff"), // Status 3 (Frozen): desaturated/icy versions

		// AntHill (15)
		ColorFromHex("#160c14ff"), ColorFromHex("#1c0815ff"), ColorFromHex("#27091dff"), ColorFromHex("#2a0920ff"),
		ColorFromHex("#140e09ff"), ColorFromHex("#1d140cff"), ColorFromHex("#20170fff"), ColorFromHex("#271d14ff"),
		ColorFromHex("#140e09ff"), ColorFromHex("#1d140cff"), ColorFromHex("#20170fff"), ColorFromHex("#271d14ff"),
		ColorFromHex("#140e09ff"), ColorFromHex("#1d140cff"), ColorFromHex("#20170fff"), ColorFromHex("#271d14ff"),
	}
)
