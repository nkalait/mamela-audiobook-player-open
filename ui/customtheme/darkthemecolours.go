package customtheme

import (
	"image/color"
)

var (
	colourBlack        = color.RGBA{R: 33, G: 33, B: 33, A: 255}
	colourBlackLight   = color.RGBA{R: 38, G: 38, B: 38, A: 255}
	colourBlackLighter = color.RGBA{R: 51, G: 51, B: 51, A: 255}
	colourWhite        = color.RGBA{R: 205, G: 205, B: 205, A: 255}
)

var (
	darkThemeColourBg      = colourBlack
	darkThemeColourFg      = colourWhite
	darkThemeColourBgLight = colourBlackLight
	// darkThemeColourBgLight = color.RGBA{244, 107, 221, 255} //colourBlackLight
	darkThemeColourBgLighter = colourBlackLighter
	darkThemeColourText      = colourWhite
)

const (
	ColourNameBackgroundLight = "bgLight"
)
