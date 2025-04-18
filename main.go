package main

import (
	"GoMaze/maze"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	
	a := app.New()
	w := a.NewWindow("Hello")
	// c := w.Canvas()

	// img := image.NewRGBA(image.Rect(0, 0, 600, 400))
	// canvasImg := canvas.NewImageFromImage(img)
	// canvasImg.FillMode = canvas.ImageFillOriginal


	maze := maze.NewWallExtending(31, 31)
	maze.Create()

	img := maze.Draw()
	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillOriginal
	w.SetContent(
		canvasImg,
	)
	w.ShowAndRun()
}