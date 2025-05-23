package maze

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)


type Follower struct {
	Rows  int
	Cols  int
	Grid [][]int
}

func NewFollower(rows int, cols int) *Follower {
	maze := &Follower{
		Rows: rows,
		Cols: cols,
	}
	maze.initGrid()
	return maze
}

func (f *Follower) initGrid() {
	grid := make([][]int, f.Rows)

	for r := range grid {
		grid[r] = make([]int, f.Cols)
		grid[r][0] = WALL
		grid[r][f.Cols-1] = WALL

		if r == 0 || r == f.Rows-1 {
			for c := 1; c <= f.Cols-2; c++ {
				grid[r][c] = WALL
			}
		}
	}
	f.Grid = grid
}


func (f *Follower) replaceValue(minR, maxR, minC, maxC, v int) {
	for j := minR; j < maxR + 1; j++ {
		for i := minC; i < maxC + 1; i++ {
			if f.Grid[j][i] == EXTENDING {
				f.Grid[j][i] = v
			}
		}
	}
}


func (f *Follower) followWall(orgR int, orgC int) {
	r, c := orgR, orgC
	minR, minC := orgR, orgC
	maxR, maxC := orgR, orgC

	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	idxes := [4]int{0, 0, 0, 0}
	var idx int

	for {
		switch f.Grid[r][c] {
		case PASSAGE:
			f.Grid[r][c] = EXTENDING

		case WALL:
			f.replaceValue(minR, maxR, minC, maxC, WALL)
			return
		}

		switch cnt := f.findExtendableDirs(&dirs, &idxes, r, c); cnt {
		case 0:
			f.replaceValue(minR, maxR, minC, maxC, PASSAGE)
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
		f.Grid[r + dr][c + dc] = EXTENDING
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

func (f *Follower) findExtendableDirs(dirs *[4][2]int, idxes *[4]int, r int, c int) int {
	cnt := 0
	
	for i, dir := range(dirs) {
		nr := r + dir[0] * 2
		nc := c + dir[1] * 2
		
		if f.Grid[nr][nc] != EXTENDING {
			idxes[cnt] = i
			cnt ++
		}
	}
	
	return cnt
}

func (f *Follower) Create() *image.RGBA {
	pts := GetStartPts(f.Rows, f.Cols)
	rand.Shuffle(len(pts), func(i, j int) {pts[i], pts[j] = pts[j], pts[i]})

	for _, pt := range(pts) {
		r, c := pt[0], pt[1]
		if f.Grid[r][c] != WALL {
			f.followWall(r, c)
		}
	}

	img := f.draw()
	return img
}


func (f *Follower) draw() *image.RGBA {
	green := color.NRGBA{R: 0, G: 128, B: 0, A: 255}
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	img := image.NewRGBA(image.Rect(0, 0, RECSIZE * f.Rows, RECSIZE * f.Cols))
	var color color.NRGBA

	for r := 0; r <= f.Rows - 1; r++ {
		for c := 0; c <= f.Cols - 1; c++ {	
			switch {
			case r == 0 && c == 1 || r == f.Rows - 1 && c == f.Cols - 2:
				color = green

			case f.Grid[r][c] == WALL:
				color = black

			default:
				color = green
			}

			draw.Draw(
				img,
				image.Rect(RECSIZE * c, RECSIZE * r, RECSIZE * c + RECSIZE, RECSIZE * r + RECSIZE),
				&image.Uniform{color},
				image.Point{},
				draw.Src,
			)
		}
	}

	return img
}