package customtheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	fyneTheme "fyne.io/fyne/v2/theme"
)

type darkTheme struct{}

func (t darkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case fyneTheme.ColorNameBackground:
		return backgroundColor()
	case fyneTheme.ColorNameForeground:
		return foregroundColor()
	case ColourNameBackgroundLight:
		return backgroundLightColor()
	case fyneTheme.ColorNameButton:
		return buttonColor()
	}
	return fyneTheme.DefaultTheme().Color(name, variant)
}

func (t darkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return fyneTheme.DefaultTheme().Font(style)
}

func (t darkTheme) Size(name fyne.ThemeSizeName) float32 {
	return fyneTheme.DefaultTheme().Size(name)
}

func (t darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// if name == theme.IconNameHome {
	// 	return fyne.NewStaticResource("myHome", homeBytes)
	// }

	return fyneTheme.DefaultTheme().Icon(name)
}

func DarkTheme() fyne.Theme {
	return &darkTheme{}
}
