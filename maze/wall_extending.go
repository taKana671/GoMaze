package maze

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

const (
	Wall      = 1
	Passage   = 0
	Extending = 2
)

type Maze interface {
	Create()
	Draw() *image.RGBA
}

type WallExtending struct {
	Rows  int
	Cols  int
	Grid [][]int
}

func NewWallExtending(rows int, cols int) *WallExtending {
	maze := &WallExtending{
		Rows: rows,
		Cols: cols,
	}
	maze.initGrid()
	return maze
}

func (w *WallExtending) initGrid() {
	grid := make([][]int, w.Rows)

	for r := range grid {
		grid[r] = make([]int, w.Cols)
		grid[r][0] = Wall
		grid[r][w.Cols-1] = Wall

		if r == 0 || r == w.Rows-1 {
			for c := 1; c <= w.Rows-2; c++ {
				grid[r][c] = 1
			}
		}
	}
	w.Grid = grid
}

func (w *WallExtending) getStartPts() [][]int {
	rows := (w.Rows - 2) / 2
	cols := (w.Cols - 2) / 2
	pts := make([][]int, rows * cols)
	i := 0

	for r := 1; r < w.Rows - 1; r++ {
		for c := 1; c < w.Cols - 1; c++ {
			if r % 2 == 0 && c % 2 == 0 {
				pts[i] = make([]int, 2)
				pts[i][0] = r
				pts[i][1] = c
				i++
			}
		} 
	}
    return pts
}


func (w *WallExtending) replaceValue(minR, maxR, minC, maxC, v int) {
	for j := minR; j < maxR + 1; j++ {
		for i := minC; i < maxC + 1; i++ {
			if w.Grid[j][i] == Extending {
				w.Grid[j][i] = v
			}
		}
	}
}


func (w *WallExtending) extendWall(orgR int, orgC int) {
	r, c := orgR, orgC
	minR, minC := orgR, orgC
	maxR, maxC := orgR, orgC

	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	idxes := [4]int{0, 0, 0, 0}
	var idx int

	for {
		switch w.Grid[r][c] {
		case Passage:
			w.Grid[r][c] = Extending

		case Wall:
			w.replaceValue(minR, maxR, minC, maxC, Wall)
			return
		}

		switch cnt := w.findExtendableDirections(&dirs, &idxes, r, c); cnt {
		case 0:
			w.replaceValue(minR, maxR, minC, maxC, Passage)
			r, c = orgR, orgC
			continue

		case 1:
			idx = idxes[cnt - 1]

		default:
			n := rand.Intn(cnt)
			idx = idxes[n]
		}

		dr := dirs[idx][0]
		dc := dirs[idx][1]
		w.Grid[r + dr][c + dc] = Extending
		r += dr * 2
		c += dc * 2

		if maxR < r {
			maxR = r
		}

		if minR > r {
			minR = r
		}

		if maxC < c {
			maxC = c
		}

		if minC > c {
			minC = c
		}

	}
}

func (w *WallExtending) findExtendableDirections(dirs *[4][2]int, idxes *[4]int, r int, c int) int {
	cnt := 0
	
	for i, dir := range(dirs) {
		nr := r + dir[0] * 2
		nc := c + dir[1] * 2
		
		if w.Grid[nr][nc] != Extending {
			idxes[cnt] = i
			cnt ++
		}
	}
	
	return cnt
}

func (w *WallExtending) Create() {
	pts := w.getStartPts()
	rand.Shuffle(len(pts), func(i, j int) {pts[i], pts[j] = pts[j], pts[i]})

	for _, pt := range(pts) {
		r, c := pt[0], pt[1]
		if w.Grid[r][c] != Wall {
			w.extendWall(r, c)
		}
	}
}

// func (w *WallExtending) Draw(img *image.RGBA) {
// 		green := color.RGBA{R: 0, G: 128, B: 0, A: 255}
// 		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
// 		side := 20

// 		var cl color.RGBA

// 		for r := 0; r <= w.Rows - 1; r++ {
// 			for c := 0; c <= w.Cols - 1; c++ {
// 				switch {
// 				case r == 0 && c == 1 || r == w.Rows - 1 && c == w.Cols - 2:
// 					cl = green
	
// 				case w.Grid[r][c] == Wall:
// 					cl = black
	
// 				default:
// 					cl= green
// 				}

// 				for y := side * r; r < side * r + side; y++ {
// 					for x := side * c; x < side * c + side; x++ {
// 						img.Set(x, y, cl)
// 					}
// 				}
// 			}
// 		}
// }

func (w *WallExtending) Draw() *image.RGBA {
	green := color.NRGBA{R: 0, G: 128, B: 0, A: 255}
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	size := 15
	img := image.NewRGBA(image.Rect(0, 0, size * w.Rows, size * w.Cols))
	var color color.NRGBA

	for r := 0; r <= w.Rows - 1; r++ {
		for c := 0; c <= w.Cols - 1; c++ {	
			switch {
			case r == 0 && c == 1 || r == w.Rows - 1 && c == w.Cols - 2:
				color = green

			case w.Grid[r][c] == Wall:
				color = black

			default:
				color = green
			}

			draw.Draw(
				img,
				image.Rect(size * c, size * r, size * c + size, size * r + size),
				&image.Uniform{color},
				image.Point{},
				draw.Src,
			)


			// draw.Draw(
			// 	img,
			// 	image.Rect(size * c, size * r, size * c + size, size * r + size),
			// 	&image.Uniform{color},
			// 	image.Point{},
			// 	draw.Src,
			// )
		}
	}

	return img
}