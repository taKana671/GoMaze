package maze

import (
	"fmt"
	"math/rand"
)

const (
	Wall      = 1
	Passage   = 0
	Extending = 2
)

type Grid struct {
	Rows  int
	Cols  int
	Array [][]int
}

func newGrid(rows int, cols int) *Grid {
	grid := &Grid{
		Rows: rows,
		Cols: cols,
	}
	grid.initArray()

	return grid
}

func (g *Grid) initArray() {
	array := make([][]int, g.Rows)

	for r := range array {
		array[r] = make([]int, g.Cols)
		array[r][0] = Wall
		array[r][g.Cols-1] = Wall

		if r == 0 || r == g.Rows-1 {
			for c := 1; c <= g.Rows-2; c++ {
				array[r][c] = 1
			}
		}
	}
	g.Array = array
}

func (g *Grid) getStartPts() [][]int {
	rows := (g.Rows - 2) / 2
	cols := (g.Cols - 2) / 2
	pts := make([][]int, rows * cols)
	i := 0

	for r := 1; r < g.Rows - 1; r++ {
		for c := 1; c < g.Cols - 1; c++ {
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


func (g *Grid) replaceValue(min_r, max_r, min_c, max_c, v int) {
	for j := min_r; j < max_r + 1; j++ {
		for i := min_c; i < max_c + 1; i++ {
			if v == Extending {
				g.Array[j][i] = v
			}
		}
	}
}


func (g *Grid) extendWall(org_r int, org_c int) {
	r, c := org_r, org_c
	min_r, min_c := org_r, org_c
	max_r, max_c := org_r, org_c

	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	idxes := [4]int{0, 0, 0, 0}
	var idx int

	for {
		switch g.Array[r][c] {
		case Passage:
			g.Array[r][c] = Extending

		case Wall:
			g.replaceValue(min_r, max_r, min_c, max_c, Wall)
			return
		}

		switch cnt := g.findExtendableDirections(&dirs, &idxes, r, c); cnt {
		case 0:
			g.replaceValue(min_r, max_r, min_c, max_c, Passage)
			r, c = org_r, org_c
			continue

		case 1:
			idx = idxes[cnt - 1]

		default:
			n := rand.Intn(cnt)
			idx = idxes[n]
		}

		dr := dirs[idx][0]
		dc := dirs[idx][1]
		g.Array[r + dr][c + dc] = Extending
		r += dr * 2
		c += dc * 2

		if max_r < r {
			max_r = r
		}

		if min_r > r {
			min_r = r
		}

		if max_c < c {
			max_c = c
		}

		if min_c > c {
			min_c = c
		}

	}
}

func (g *Grid) findExtendableDirections(dirs *[4][2]int, idxes *[4]int, r int, c int) int {
	cnt := 0
	
	for i, dir := range(dirs) {
		nr := r + dir[0] * 2
		nc := c + dir[1] * 2
		
		if g.Array[nr][nc] != Extending {
			idxes[cnt] = i
			cnt ++
		}
	}
	
	return cnt
}

func CreateMaze(rows int, cols int) {
	grid := newGrid(rows, cols)
	pts := grid.getStartPts()
	fmt.Println(grid)
	fmt.Println(pts)

	rand.Shuffle(len(pts), func(i, j int) {pts[i], pts[j] = pts[j], pts[i]})
	fmt.Println(pts)

	for _, pt := range(pts) {
		r, c := pt[0], pt[1]
		if grid.Array[r][c] != Wall {
			grid.extendWall(r, c)
		}
	}
}