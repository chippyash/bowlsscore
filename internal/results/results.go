/**
Bowls Scorer
copyright (c) Ashley Kitson, UK, 2023
Licence: BSD-3-Clause See LICENSE
*/

package results

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/chippyash/bowlsscore/internal/constants"
	"github.com/chippyash/bowlsscore/internal/scoring"
	"strconv"
)

type results struct {
	app       fyne.App
	win       fyne.Window
	container *fyne.Container
	scores    *scoring.Scores
}

func New(app fyne.App, win fyne.Window, c *fyne.Container, scores *scoring.Scores) *results {
	return &results{app: app, win: win, container: c, scores: scores}
}

func (r *results) calcGame() [][3]string {
	var data [100][3]string
	s := [3]string{constants.LblEnds, constants.LblHome, constants.LblAway}
	data[0] = s
	var prevH, prevA int
	for x := 1; x <= len(*r.scores); x++ {
		v := (*r.scores)[x]
		thisHome := v.Home - prevH
		thisAway := v.Away - prevA
		prevH = v.Home
		prevA = v.Away
		strH := fmt.Sprintf("- / %d", v.Home)
		strA := fmt.Sprintf("- / %d", v.Away)
		if thisHome > 0 {
			strH = fmt.Sprintf("%d / %d", thisHome, v.Home)
		}
		if thisAway > 0 {
			strA = fmt.Sprintf("%d / %d", thisAway, v.Away)
		}
		data[x] = [3]string{strconv.Itoa(x), strH, strA}
	}
	l := len(*r.scores) + 1
	return data[:l]
}

func (r *results) calcTrack() [][3]string {
	var data [100][3]string
	s := [3]string{constants.LblEnds, constants.LblClose, constants.LblScored}
	data[0] = s
	for x := 1; x <= len(*r.scores); x++ {
		v := (*r.scores)[x]
		data[x] = [3]string{strconv.Itoa(x), strconv.Itoa(v.Home), strconv.Itoa(v.Away)}
	}
	l := len(*r.scores) + 1
	return data[:l]
}

func (r *results) calc() [][3]string {
	if r.app.Preferences().Int(constants.ModeName) == constants.ModeGame {
		return r.calcGame()
	}
	return r.calcTrack()
}

func (r *results) Display() {
	var modal *widget.PopUp
	data := r.calc()
	tble := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Results")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(data[id.Row][id.Col])
			object.(*widget.Label).Alignment = fyne.TextAlignTrailing
		})
	closeBtn := widget.NewButton("Close", func() {
		if modal != nil {
			modal.Hide()
		}
	})
	c := container.NewVSplit(
		tble,
		closeBtn,
	)
	c.SetOffset(1.0)
	modal = widget.NewModalPopUp(
		c,
		r.win.Canvas(),
	)
	modal.Resize(r.container.Size())

	modal.Show()
}
