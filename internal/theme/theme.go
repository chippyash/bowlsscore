//go:generate fyne bundle -o bundled.go bowls.png

/**
Bowls Scorer
copyright (c) Ashley Kitson, UK, 2023
Licence: BSD-3-Clause See LICENSE
*/

package theme

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

//go:embed bowls.png
var icon []byte

type BowlsTheme struct {
	fyne.Theme
	variant int //default will be zero - dark theme
}

//New returns the theme for the app
func New() *BowlsTheme {
	return &BowlsTheme{}
}

//Color implements fyne.Theme interface
func (t BowlsTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	//switch name {
	//case theme.ColorNameBackground:
	//	return color.RGBA{
	//		R: 28,
	//		G: 57,
	//		B: 115,
	//		A: 0x0000,
	//	}
	//case theme.ColorNameButton:
	//	return color.RGBA{
	//		R: 22,
	//		G: 45,
	//		B: 91,
	//		A: 0x00ff,
	//	}
	//case theme.ColorNameHover:
	//	return color.RGBA{
	//		R: 51,
	//		G: 104,
	//		B: 209,
	//		A: 0x00ff,
	//	}
	//case theme.ColorNamePressed:
	//	return color.RGBA{
	//		R: 54,
	//		G: 209,
	//		B: 209,
	//		A: 0x00ff,
	//	}
	//default:
	//	return theme.DefaultTheme().Color(name, variant)
	//}
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(t.variant))
}

func (t *BowlsTheme) SetVariant(v int) *BowlsTheme {
	t.variant = v
	return t
}

//Icon implements fyne.Theme interface
func (t BowlsTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if name == theme.IconNameHome {
		return fyne.NewStaticResource("bowls.png", icon)
	}

	return theme.DefaultTheme().Icon(name)
}

//Font implements fyne.Theme interface
func (t BowlsTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

//Size implements fyne.Theme interface
func (t BowlsTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 24
	}
	return theme.DefaultTheme().Size(name)
}

//ScoreLabel returns a label containing a score
func (t BowlsTheme) ScoreLabel(bind binding.Int) *widget.Label {
	l := widget.NewLabelWithData(binding.IntToStringWithFormat(bind, "%02d"))
	l.Alignment = fyne.TextAlignCenter
	l.TextStyle = fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
		Symbol:    false,
		TabWidth:  0,
	}

	return l
}

//ScoreBtn returns a button to manipulate a score
func (t BowlsTheme) ScoreBtn(label string, fn func(), icon fyne.Resource) *widget.Button {
	b := widget.NewButtonWithIcon(label, icon, fn)
	b.Importance = widget.MediumImportance
	return b
}
