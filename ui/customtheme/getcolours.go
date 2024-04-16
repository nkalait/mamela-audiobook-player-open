package customtheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	fyneTheme "fyne.io/fyne/v2/theme"
)

func GetColour(colorName fyne.ThemeColorName) color.Color {
	if IsDark() {
		return darkTheme{}.Color(colorName, fyneTheme.VariantDark)
	}
	if IsLight() {
		return darkTheme{}.Color(colorName, fyneTheme.VariantLight)
	}
	return darkTheme{}.Color(colorName, fyneTheme.VariantDark)
}

func IsDark() bool {
	return fyne.CurrentApp().Settings().ThemeVariant() == 0
}

func IsLight() bool {
	return fyne.CurrentApp().Settings().ThemeVariant() == 1

}

func backgroundColor() color.Color {
	if IsDark() {
		return darkThemeColourBg
	}
	if IsLight() {
		return lightThemeColourBg
	}
	return fyneTheme.BackgroundColor()
}

func backgroundLightColor() color.Color {
	if IsDark() {
		return darkThemeColourBgLight
	}
	if IsLight() {
		return lightThemeColourBgLight
	}
	return fyneTheme.BackgroundColor()
}

func foregroundColor() color.Color {
	if IsDark() {
		return darkThemeColourFg
	}
	if IsLight() {
		return lightThemeColourFg
	}
	return fyneTheme.ForegroundColor()
}

func buttonColor() color.Color {
	if IsDark() {
		return darkThemeColourButton
	}
	if IsLight() {
		return lightThemeColourButton
	}
	return fyneTheme.ForegroundColor()
}
