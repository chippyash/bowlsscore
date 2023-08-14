/**
Bowls Scorer
copyright (c) Ashley Kitson, UK, 2023
Licence: BSD-3-Clause See LICENSE
*/

package prefs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	theme2 "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/chippyash/bowlsscore/internal/constants"
	"github.com/chippyash/bowlsscore/internal/scorecard"
	"github.com/chippyash/bowlsscore/internal/theme"
)

type Prefs struct {
	win       fyne.Window
	app       fyne.App
	scorecard *scorecard.ScoreCard
}

func New(app fyne.App, window fyne.Window) *Prefs {
	p := &Prefs{
		win: window,
		app: app,
	}
	//set defaults
	p.setBool(constants.AutoSaveName, p.asBoolWithFallback(constants.AutoSaveName, false))
	p.setInt(constants.AutoSecsName, p.asIntWithFallback(constants.AutoSecsName, 5))
	p.setInt(constants.VariantName, p.asIntWithFallback(constants.VariantName, int(theme2.VariantDark)))
	p.setInt(constants.ModeName, p.asIntWithFallback(constants.ModeName, constants.ModeGame))

	return p
}

func (p *Prefs) WithScorecard(scorecard *scorecard.ScoreCard) *Prefs {
	p.scorecard = scorecard
	return p
}

func (p *Prefs) AutoSave() bool {
	return p.app.Preferences().Bool(constants.AutoSaveName)
}

func (p *Prefs) AutoSecs() int {
	return p.app.Preferences().Int(constants.AutoSecsName)
}

func (p *Prefs) ThemeVariant() int {
	return p.app.Preferences().Int(constants.VariantName)
}

func (p *Prefs) PlayMode() int {
	return p.app.Preferences().Int(constants.ModeName)
}

func (p *Prefs) asBoolWithFallback(name string, fallback bool) bool {
	return p.app.Preferences().BoolWithFallback(name, fallback)
}

func (p *Prefs) asIntWithFallback(name string, fallback int) int {
	return p.app.Preferences().IntWithFallback(name, fallback)
}

func (p *Prefs) setBool(name string, value bool) {
	p.app.Preferences().SetBool(name, value)
}

func (p *Prefs) setInt(name string, value int) {
	p.app.Preferences().SetInt(name, value)
}

func (p *Prefs) Display() {
	items := make([]*widget.FormItem, 0)
	saveState := p.AutoSave()
	autoSaveWidget := widget.NewRadioGroup([]string{"Yes", "No"}, func(opt string) {
		saveState = opt == "Yes"
	})
	if saveState {
		autoSaveWidget.SetSelected("Yes")
	} else {
		autoSaveWidget.SetSelected("No")
	}

	items = append(items, widget.NewFormItem("Auto Save", autoSaveWidget))

	autoSecsLable := binding.NewString()
	autoSecsLableWidget := widget.NewLabelWithData(autoSecsLable)
	_ = autoSecsLable.Set(fmt.Sprintf("%d secs", p.AutoSecs()))
	autoSecs := binding.NewFloat()
	_ = autoSecs.Set(float64(p.AutoSecs()))
	autoSecsWidget := widget.NewSliderWithData(1.0, 15.0, autoSecs)
	autoSecsWidget.OnChanged = func(v float64) {
		_ = autoSecs.Set(v)
		_ = autoSecsLable.Set(fmt.Sprintf("%d secs", int(v)))
	}

	variantState := p.ThemeVariant()
	variantWidget := widget.NewRadioGroup([]string{"Dark", "Light"}, func(opt string) {
		if opt == "Dark" {
			variantState = int(theme2.VariantDark)
			return
		}
		variantState = int(theme2.VariantLight)
	})
	if variantState == int(theme2.VariantDark) {
		variantWidget.SetSelected("Dark")
	} else {
		variantWidget.SetSelected("Light")
	}

	playMode := p.PlayMode()
	playWidget := widget.NewRadioGroup([]string{"Game", "Track"}, func(opt string) {
		if opt == "Game" {
			playMode = constants.ModeGame
			return
		}
		playMode = constants.ModeTrack
	})
	if playMode == constants.ModeGame {
		playWidget.SetSelected("Game")
	} else {
		playWidget.SetSelected("Track")
	}

	items = append(items, widget.NewFormItem("Auto Save Seconds", autoSecsLableWidget))
	items = append(items, widget.NewFormItem("", autoSecsWidget))
	items = append(items, widget.NewFormItem("Theme", variantWidget))
	items = append(items, widget.NewFormItem("Play Mode", playWidget))

	var modal *widget.PopUp
	form := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items:      items,
		OnSubmit: func() {
			//set the preferences
			asecs, _ := autoSecs.Get()
			p.setInt(constants.AutoSecsName, int(asecs))
			p.setBool(constants.AutoSaveName, saveState)
			p.setInt(constants.VariantName, variantState)
			p.setInt(constants.ModeName, playMode)

			//hide the dialog
			modal.Hide()
			//reset the scorecard display
			scores := p.scorecard.GetScores()
			t := p.app.Settings().Theme().(*theme.BowlsTheme).SetVariant(variantState)

			scoreCard := scorecard.New(t, p.app).SetScores(scores)
			layout := scoreCard.Layout()
			p.scorecard = scoreCard
			//refresh the display
			p.win.SetContent(layout)
			p.win.Canvas().Refresh(p.win.Canvas().Content())
		},
		OnCancel:   func() { modal.Hide() },
		SubmitText: "Save",
		CancelText: "Cancel",
	}

	modal = widget.NewModalPopUp(
		form,
		p.win.Canvas(),
	)
	modal.Resize(p.win.Content().Size())
	modal.Show()
}
