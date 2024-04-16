package customtheme

import (
	"image/color"
)

var (
	colourBlack      = color.RGBA{R: 50, G: 50, B: 50, A: 255}
	colourBlackLight = color.RGBA{R: 58, G: 58, B: 58, A: 255}
)
var (
	colourWhite      = color.RGBA{R: 233, G: 233, B: 233, A: 255}
	colourWhiteLight = color.RGBA{R: 225, G: 225, B: 225, A: 255}
)
var (
	darkThemeColourBg      = colourBlack
	darkThemeColourFg      = colourWhite
	darkThemeColourBgLight = colourBlackLight
	darkThemeColourButton  = colourBlack
)
var (
	lightThemeColourBg      = colourWhite
	lightThemeColourBgLight = colourWhiteLight
	lightThemeColourFg      = colourBlack
	lightThemeColourButton  = colourWhite
)

const (
	ColourNameBackgroundLight = "bgLight"
)
