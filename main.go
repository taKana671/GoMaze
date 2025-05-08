package main

import (
	"GoMaze/maize"
	"image"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)


const (
	FOLLOWING = "Wall Following"
	DIGGING = "Wall Digging"
)


const (
	ROWS = 31
	COLS = 31
)


func main() {
	a := app.New()
	w := a.NewWindow("GoMaze")

	img := image.NewRGBA(image.Rect(0, 0,  maize.RECSIZE * ROWS, maize.RECSIZE * COLS))
	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillOriginal

	genMaze := func(algo string) {
		var img *image.RGBA
		
		switch algo {
		case FOLLOWING:
			img = maize.NewFollower(31, 31).Create()
		case DIGGING:
			img = maize.NewDigger(31, 31).Create()
		}

		canvasImg.Image = img
		canvas.Refresh(canvasImg)
	}

	sl := widget.NewSelect([]string{
		FOLLOWING,
		DIGGING,
		}, func(s string) {
			genMaze(s)
		})

	sl.SetSelected(FOLLOWING)

	tb := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			genMaze(sl.Selected)
		}),
	)

	mainContainer := container.New(
		layout.NewBorderLayout(
			sl, tb, nil, nil,
		),
		sl,
		canvasImg,
		tb,
	)

	w.SetContent(mainContainer)
	w.ShowAndRun()
}