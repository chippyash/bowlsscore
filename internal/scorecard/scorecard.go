/**
Bowls Scorer
copyright (c) Ashley Kitson, UK, 2023
Licence: BSD-3-Clause See LICENSE
*/

package scorecard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/chippyash/bowlsscore/internal/constants"
	"github.com/chippyash/bowlsscore/internal/scoring"
	btheme "github.com/chippyash/bowlsscore/internal/theme"
	"time"
)

const (
	btnHomeInc = iota
	btnHomeDec
	btnAwayInc
	btnAwayDec
	btnEndsInc
	btnEndsDec
)

type buttons [6]*widget.Button
type display [3]*widget.Label

type ScoreCard struct {
	home    binding.Int
	ends    binding.Int
	away    binding.Int
	theme   *btheme.BowlsTheme
	app     fyne.App
	scores  scoring.Scores
	buttons buttons
	display display
}

//autosaving
var autoSave bool
var autoSecs int64

//New returns a pointer to a new scoreCard
func New(t *btheme.BowlsTheme, app fyne.App) *ScoreCard {
	s := &ScoreCard{
		home:   binding.NewInt(),
		ends:   binding.NewInt(),
		away:   binding.NewInt(),
		theme:  t,
		scores: make(scoring.Scores),
		app:    app,
	}
	_ = s.home.Set(0)
	_ = s.ends.Set(0)
	_ = s.away.Set(0)

	return s
}

//incScore increments a score with 99 as an upper bound
func (s *ScoreCard) incScore(score binding.Int, btn int) {
	v, _ := score.Get()
	v++
	if v > 99 {
		v = 99
	}
	_ = score.Set(v)

	if autoSave {
		s.AutoSave()
		//s.buttons[btn].Importance = widget.WarningImportance
		//s.buttons[btn].Refresh()
	}
}

//decScore decrements a score with zero as a lower bound
func (s *ScoreCard) decScore(score binding.Int, btn int) {
	v, _ := score.Get()
	v--
	if v < 0 {
		v = 0
	}
	_ = score.Set(v)
	if autoSave {
		s.AutoSave()
	}
}

//Reset resets all scores to zero
func (s *ScoreCard) Reset() {
	_ = s.home.Set(0)
	_ = s.ends.Set(0)
	_ = s.away.Set(0)
	s.scores = make(scoring.Scores)
}

func (s *ScoreCard) SaveScore() {
	e, _ := s.ends.Get()
	h, _ := s.home.Get()
	a, _ := s.away.Get()
	s.scores[e] = scoring.Score{
		Home: h,
		Away: a,
	}
	//if autoSave {
	//	for b := btnHomeInc; b <= btnAwayDec; b++ {
	//		s.buttons[b].Importance = widget.MediumImportance
	//		s.buttons[b].Refresh()
	//	}
	//}
}

var autoSaveAlreadyRunning bool

//AutoSave sets a timer to automatically save the scores
func (s *ScoreCard) AutoSave() {
	if autoSaveAlreadyRunning {
		return
	}
	autoSaveAlreadyRunning = true
	go func(t time.Time) {
		for {
			now := time.Now()
			if now.After(t) {
				end, _ := s.ends.Get()
				end++
				_ = s.ends.Set(end)
				s.SaveScore()
				autoSaveAlreadyRunning = false
				break
			}
			time.Sleep(time.Millisecond * 5)
		}
	}(time.Now().Add(time.Second * time.Duration(autoSecs)))
}

//Layout returns the container for the application
func (s *ScoreCard) Layout() *fyne.Container {
	//preferences
	autoSave = s.app.Preferences().Bool(constants.AutoSaveName)
	autoSecs = int64(s.app.Preferences().Int(constants.AutoSecsName))
	playMode := s.app.Preferences().Int(constants.ModeName)

	//buttons
	s.buttons[btnHomeInc] = s.theme.ScoreBtn("", func() { s.incScore(s.home, btnHomeInc) }, theme.ContentAddIcon())
	s.buttons[btnHomeDec] = s.theme.ScoreBtn("", func() { s.decScore(s.home, btnHomeDec) }, theme.ContentRemoveIcon())
	s.buttons[btnEndsInc] = s.theme.ScoreBtn("", func() { s.incScore(s.ends, constants.DispEnds) }, theme.ContentAddIcon())
	s.buttons[btnEndsDec] = s.theme.ScoreBtn("", func() { s.decScore(s.ends, constants.DispEnds) }, theme.ContentRemoveIcon())
	if autoSave {
		s.buttons[btnEndsDec] = s.theme.ScoreBtn("", func() {}, nil)
		s.buttons[btnEndsDec].Disable()
		s.buttons[btnEndsInc] = s.theme.ScoreBtn("", func() {}, nil)
		s.buttons[btnEndsInc].Disable()
	}
	s.buttons[btnAwayInc] = s.theme.ScoreBtn("", func() { s.incScore(s.away, constants.DispAway) }, theme.ContentAddIcon())
	s.buttons[btnAwayDec] = s.theme.ScoreBtn("", func() { s.decScore(s.away, constants.DispAway) }, theme.ContentRemoveIcon())
	btnReset := s.theme.ScoreBtn("", func() { s.Reset() }, theme.DeleteIcon())
	btnSave := s.theme.ScoreBtn("", func() { s.SaveScore() }, theme.DocumentSaveIcon())
	if autoSave {
		btnSave.Hide()
	}

	//data binding
	s.display[constants.DispHome] = s.theme.ScoreLabel(s.home)
	s.display[constants.DispEnds] = s.theme.ScoreLabel(s.ends)
	s.display[constants.DispAway] = s.theme.ScoreLabel(s.away)

	var hTitle, aTitle string
	if playMode == constants.ModeGame {
		hTitle = constants.LblHome
		aTitle = constants.LblAway
	} else {
		hTitle = constants.LblClose
		aTitle = constants.LblScored
	}
	homeTitle := container.NewGridWithColumns(
		1,
		widget.NewLabelWithStyle(hTitle, fyne.TextAlignCenter, fyne.TextStyle{
			Bold:      true,
			Italic:    false,
			Monospace: true,
			Symbol:    false,
			TabWidth:  0,
		}),
	)
	homeScores := container.NewGridWithColumns(
		3,
		s.buttons[btnHomeDec], s.display[constants.DispHome], s.buttons[btnHomeInc],
	)
	endsTitle := container.NewGridWithColumns(
		1,
		widget.NewLabelWithStyle(constants.LblEnds, fyne.TextAlignCenter, fyne.TextStyle{
			Bold:      false,
			Italic:    false,
			Monospace: true,
			Symbol:    false,
			TabWidth:  0,
		}),
	)
	endsScores := container.NewGridWithColumns(
		3,
		s.buttons[btnEndsDec], s.display[constants.DispEnds], s.buttons[btnEndsInc],
	)
	awayTitle := container.NewGridWithColumns(
		1,
		widget.NewLabelWithStyle(aTitle, fyne.TextAlignCenter, fyne.TextStyle{
			Bold:      true,
			Italic:    false,
			Monospace: true,
			Symbol:    false,
			TabWidth:  0,
		}),
	)
	awayScores := container.NewGridWithColumns(
		3,
		s.buttons[btnAwayDec], s.display[constants.DispAway], s.buttons[btnAwayInc],
	)
	bottom := container.NewGridWithColumns(2, btnReset, btnSave)

	//the main layout
	l := layout.NewGridLayoutWithRows(7)

	return container.New(l, homeTitle, homeScores, endsTitle, endsScores, awayTitle, awayScores, bottom)
}

func (s *ScoreCard) GetScores() scoring.Scores {
	return s.scores
}

func (s *ScoreCard) SetScores(scores scoring.Scores) *ScoreCard {
	s.scores = scores
	e := len(scores)
	_ = s.ends.Set(e)
	if e > 0 {
		_ = s.home.Set(scores[e].Home)
		_ = s.away.Set(scores[e].Away)
	}
	return s
}

func (s *ScoreCard) GetEnds() int {
	ends, _ := s.ends.Get()
	return ends
}

func (s *ScoreCard) SetEnds(ends int) *ScoreCard {
	_ = s.ends.Set(ends)
	s.display[constants.DispEnds] = s.theme.ScoreLabel(s.ends)
	return s
}
