package maze

import (
	"math/rand"
	"time"
)

type Point struct {
	X, Y int
}

type Maze struct {
	Width    int
	Height   int
	Grid     [][]bool
	Position Point
	Path     []Point
	rng      *rand.Rand
}

func NewMaze(initialPosition Point, screenWidth, blockSize int) *Maze {
	width := screenWidth / blockSize

	grid := make([][]bool, width)

	for i := range grid {
		grid[i] = make([]bool, width)

		for j := range width {
			grid[i][j] = false
		}
	}

	// creates a seed (source) with the source being the current time in nanoseconds
	s := rand.NewSource(time.Now().UnixNano())
	// creates instance of rand with the designated seed
	rng := rand.New(s)

	return &Maze{
		Width:    width,
		Height:   width,
		Grid:     grid,
		rng:      rng,
		Position: initialPosition,
		Path:     []Point{initialPosition},
	}
}

func (m *Maze) Carve() {
	y := m.Position.Y
	x := m.Position.X

	dirs := []Point{
		{X: 0, Y: -2}, // up
		{X: 0, Y: 2},  // down
		{X: 2, Y: 0},  // right
		{X: -2, Y: 0}, // left
	}

	// shuffles with Fisher-Yarnes algorithm
	m.rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, direction := range dirs {
		nx, ny := x+direction.X, y+direction.Y

		canGoToX := nx > 0 && nx < m.Width-1
		canGoToY := ny > 0 && ny < m.Height-1

		if canGoToX && canGoToY && !m.Grid[ny][nx] {
			intermediateY := y + direction.Y/2
			intermediateX := x + direction.X/2

			m.Grid[y][x] = true
			m.Grid[ny][nx] = true
			m.Grid[intermediateY][intermediateX] = true

			m.Position.X = nx
			m.Position.Y = ny

			m.Path = append(m.Path, Point{X: nx, Y: ny})

			return
		}
	}

	m.Backtrack()
}

func (m *Maze) Backtrack() {
	if len(m.Path) <= 1 {
		m.Grid[1][0] = true
		m.Grid[len(m.Grid)-2][len(m.Grid)-1] = true
		return
	}

	m.Path = m.Path[:len(m.Path)-1]

	lastPosition := m.Path[len(m.Path)-1]
	m.Position = lastPosition
}
