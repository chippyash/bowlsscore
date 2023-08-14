/**
Bowls Scorer
As simple lawn green bowls score recorder

copyright (c) Ashley Kitson, UK, 2023
Licence: BSD-3-Clause See LICENSE
*/

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/chippyash/bowlsscore/internal/prefs"
	"github.com/chippyash/bowlsscore/internal/results"
	"github.com/chippyash/bowlsscore/internal/scorecard"
	"github.com/chippyash/bowlsscore/internal/theme"
)

func main() {
	//a := app.New()
	a := app.NewWithID("biz.zf4.bowlsscore")
	t := theme.New()
	a.Settings().SetTheme(t)
	w := a.NewWindow(a.Metadata().Name)

	//set preferences
	prefs := prefs.New(a, w)
	t.SetVariant(prefs.ThemeVariant())

	//create main screen and fetch layout
	scoreCard := scorecard.New(t, a)
	prefs = prefs.WithScorecard(scoreCard)
	layout := scoreCard.Layout()

	mm := fyne.NewMainMenu(
		fyne.NewMenu("Menu",
			fyne.NewMenuItem("Results", func() {
				scores := scoreCard.GetScores()
				res := results.New(a, w, layout, &scores)
				res.Display()
			}),
			fyne.NewMenuItem("Preferences", func() {
				prefs.Display()
			}),
		),
	)
	w.SetMainMenu(mm)

	//resize if displaying on PC screen - useful when developing
	if !a.Driver().Device().IsMobile() && !a.Driver().Device().IsBrowser() {
		w.Resize(fyne.Size{
			Width:  500,
			Height: 300,
		})
	}

	w.SetContent(layout)
	w.ShowAndRun()
}
