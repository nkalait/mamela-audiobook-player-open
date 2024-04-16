package customtheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	fyneTheme "fyne.io/fyne/v2/theme"
)

type lightTheme struct{}

func (t lightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case fyneTheme.ColorNameBackground:
		return backgroundColor()
	case fyneTheme.ColorNameForeground:
		return foregroundColor()
	case ColourNameBackgroundLight:
		return backgroundLightColor()
	}
	return fyneTheme.DefaultTheme().Color(name, variant)
}

func (t lightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return fyneTheme.DefaultTheme().Font(style)
}

func (t lightTheme) Size(name fyne.ThemeSizeName) float32 {
	return fyneTheme.DefaultTheme().Size(name)
}

func (t lightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// if name == theme.IconNameHome {
	// 	return fyne.NewStaticResource("myHome", homeBytes)
	// }

	return fyneTheme.DefaultTheme().Icon(name)
}

func LightTheme() fyne.Theme {
	return &lightTheme{}
}
