package maize

import "image"


const RECSIZE = 15

const (
	WALL = 1
	PASSAGE = 0
	EXTENDING = 2
)


type Maize interface {
	Create() *image.RGBA
}


func GetStartPts(rows int, cols int) [][]int {
	rn := (rows - 2) / 2
	cn := (cols - 2) / 2
	pts := make([][]int, rn * cn)
	i := 0

	for r := 1; r < rows - 1; r++ {
		for c := 1; c < cols - 1; c++ {
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