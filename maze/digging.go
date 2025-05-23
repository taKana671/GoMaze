package maze

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"slices"
	"sync"
)


type Digger struct {
	Rows int
	Cols int
	Grid [][]int
}


const (
	up = iota
	down
	right
	left
)


type cell struct {
	x int
	y int
}


func NewDigger(rows int, cols int) *Digger {
	maze := &Digger{
		Rows: rows + 2,
		Cols: cols + 2,
	}
	maze.initGrid()
	return maze
}


func (d *Digger) initGrid() {
	grid := make([][]int, d.Rows)

	for r := range grid {
		grid[r] = make([]int, d.Cols)
		grid[r][0] = PASSAGE
		grid[r][d.Cols-1] = PASSAGE

		if r == 0 || r == d.Rows-1 {
			for c := 1; c <= d.Cols-2; c++ {
				grid[r][c] = PASSAGE
			}
		} else {
			for c := 1; c <= d.Cols-2; c++ {
				grid[r][c] = WALL
			}
		}
	}
	d.Grid = grid
}


func (d *Digger) findDiggableDirs(dirs *[4]int, c *cell) int {
	cnt := 0

	if d.Grid[c.y+1][c.x] == WALL && d.Grid[c.y+2][c.x] == WALL {
		dirs[cnt] = up
		cnt++
	}

	if d.Grid[c.y][c.x+1] == WALL && d.Grid[c.y][c.x+2] == WALL {
		dirs[cnt] = right
		cnt++
	}

	if d.Grid[c.y][c.x-1] == WALL && d.Grid[c.y][c.x-2] == WALL {
		dirs[cnt] = left
		cnt++
	}

	if d.Grid[c.y-1][c.x] == WALL && d.Grid[c.y-2][c.x] == WALL {
		dirs[cnt] = down
		cnt++
	}

	return cnt
}


func (d *Digger) digWall(startPt cell) {
	var c cell
	var i, cnt int

	d.Grid[startPt.y][startPt.x] = PASSAGE
	psgs := []cell{startPt}
	dirs := [4]int{0, 0, 0, 0}

	for {
		if len(psgs) == 0 {
			break
		}

		i = rand.Intn(len(psgs))
		c = psgs[i]
		psgs = slices.Delete(psgs, i, i + 1)

		for {
			cnt = d.findDiggableDirs(&dirs, &c)
			
			if cnt == 0 {
				break
			}

			i = rand.Intn(cnt)

			switch dirs[i] {
			case up:
				d.Grid[c.y+1][c.x] = PASSAGE
				d.Grid[c.y+2][c.x] = PASSAGE
				psgs = append(psgs, cell{c.x, c.y+2})
			case right:
				d.Grid[c.y][c.x+1] = PASSAGE
				d.Grid[c.y][c.x+2] = PASSAGE
				psgs = append(psgs, cell{c.x+2, c.y})
			case down:
				d.Grid[c.y-1][c.x] = PASSAGE
				d.Grid[c.y-2][c.x] = PASSAGE
				psgs = append(psgs, cell{c.x, c.y-2})
			case left:
				d.Grid[c.y][c.x-1] = PASSAGE
				d.Grid[c.y][c.x-2] = PASSAGE
				psgs = append(psgs, cell{c.x-2, c.y})
			}
		}
	}
}


func (d *Digger) Create() *image.RGBA {
	pts := GetStartPts(d.Rows, d.Cols)
	i := rand.Intn(len(pts))
	pt := pts[i]
	cell := cell{pt[0], pt[1]}

	d.digWall(cell)
	img := d.draw()
	return img
}

func (d *Digger) draw() *image.RGBA {
	green := color.NRGBA{R: 0, G: 128, B: 0, A: 255}
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	var wg sync.WaitGroup
	// wg.Add(31 * 31)

	img := image.NewRGBA(image.Rect(0, 0, RECSIZE * (d.Rows - 2), RECSIZE * (d.Cols - 2)))
	var color color.NRGBA
	var mu sync.Mutex


	for r := 1; r <= d.Rows - 2; r++ {
		for c := 1; c <= d.Cols - 2; c++ {
			wg.Add(1)

			switch {
			case r == 1 && c == 2 || r == d.Rows - 2 && c == d.Cols - 3:
				color = green

			case d.Grid[r][c] == WALL:
				color = black

			default:
				color = green
			}

			_c := c - 1
			_r := r - 1

			go func() {
				defer wg.Done()

				mu.Lock()
				draw.Draw(
					img,
					image.Rect(RECSIZE * _c, RECSIZE * _r, RECSIZE * _c + RECSIZE, RECSIZE * _r + RECSIZE),
					&image.Uniform{color},
					image.Point{},
					draw.Src,
				)
				mu.Unlock()
			}()
		}
	}

	wg.Wait()
	return img
}

// func (d *Digger) draw() *image.RGBA {
// 	green := color.NRGBA{R: 0, G: 128, B: 0, A: 255}
// 	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

// 	img := image.NewRGBA(image.Rect(0, 0, RECSIZE * (d.Rows - 2), RECSIZE * (d.Cols - 2)))
// 	var color color.NRGBA

// 	for r := 1; r <= d.Rows - 2; r++ {
// 		for c := 1; c <= d.Cols - 2; c++ {	
// 			switch {
// 			case r == 1 && c == 2 || r == d.Rows - 2 && c == d.Cols - 3:
// 				color = green

// 			case d.Grid[r][c] == WALL:
// 				color = black

// 			default:
// 				color = green
// 			}

// 			_c := c - 1
// 			_r := r - 1

// 			draw.Draw(
// 				img,
// 				image.Rect(RECSIZE * _c, RECSIZE * _r, RECSIZE * _c + RECSIZE, RECSIZE * _r + RECSIZE),
// 				&image.Uniform{color},
// 				image.Point{},
// 				draw.Src,
// 			)
// 		}
// 	}
// 	return img
// }