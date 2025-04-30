package main

import (
	"GoMaze/maze"
	"image"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	
	a := app.New()
	w := a.NewWindow("Hello")
	
	genMaze := func() *image.RGBA {
		grid := maze.NewWallExtending(31, 31)
		grid.Create()
		return grid.Draw()
	}

	img := genMaze()
	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillOriginal
	
	tb := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			canvasImg.Image = genMaze()
			canvas.Refresh(canvasImg)
		}),
	)
	
	mainContainer := container.New(
		layout.NewBorderLayout(
			nil, tb, nil, nil,
		),
		canvasImg,
		tb,
	)

	w.SetContent(mainContainer)
	w.ShowAndRun()
}