package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type pixTheme struct{}

func (pixTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 12, G: 121, B: 111, A: 255}
	case theme.ColorNameBackground:
		return color.NRGBA{R: 245, G: 248, B: 247, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 24, G: 44, B: 53, A: 255}
	}
	return theme.DefaultTheme().Color(n, theme.VariantLight)
}
func (pixTheme) Font(s fyne.TextStyle) fyne.Resource     { return theme.DefaultTheme().Font(s) }
func (pixTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (pixTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return 16
	}
	return theme.DefaultTheme().Size(n)
}
