package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	moleSpeed    = time.Second / 6
	screenWidth  = 300
	screenHeight = 480
	blockSize    = 20
	boardSize    = screenWidth / blockSize
	boardOffsetX = 0
	boardOffsetY = (screenHeight - screenWidth) / 2
)

type Point struct {
	x, y int
}

type Maze struct {
	width  int
	height int
	grid   [][]bool
	rng    *rand.Rand
	path   []Point
}

func NewMaze(initialPosition Point) *Maze {
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
		width:  width,
		height: width,
		grid:   grid,
		rng:    rng,
		path:   []Point{initialPosition},
	}
}

func (m *Maze) Carve(g *Game) bool {
	y := g.position.y
	x := g.position.x

	dirs := []Point{
		{0, -2}, // up
		{0, 2},  // down
		{2, 0},  // right
		{-2, 0}, // left
	}

	// shuffles with Fisher-Yarnes algorithm
	m.rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, direction := range dirs {
		nx, ny := x+direction.x, y+direction.y

		canGoToX := nx > 0 && nx < m.width-1
		canGoToY := ny > 0 && ny < m.height-1

		if canGoToX && canGoToY && !m.grid[ny][nx] {
			intermediateY := y + direction.y/2
			intermediateX := x + direction.x/2

			m.grid[y][x] = true
			m.grid[ny][nx] = true
			m.grid[intermediateY][intermediateX] = true

			g.position.x = nx
			g.position.y = ny

			g.maze.path = append(g.maze.path, Point{nx, ny})

			return true
		}
	}

	return false
}

type Game struct {
	position   Point
	lastUpdate time.Time
	maze       *Maze
}

func (g *Game) Update() error {
	if time.Since(g.lastUpdate) < moleSpeed {
		return nil
	}

	couldCarve := g.maze.Carve(g)

	if !couldCarve && len(g.maze.path) == 0 {
		return nil
	}

	if !couldCarve {
		g.maze.path = g.maze.path[:len(g.maze.path)-1]

		if len(g.maze.path) == 0 {
			return nil
		}

		lastPosition := g.maze.path[len(g.maze.path)-1]
		fmt.Println("last", lastPosition.x, lastPosition.y)
		g.position = lastPosition
	}

	g.lastUpdate = time.Now()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 49, 59, 255})

	// bg
	vector.FillRect(
		screen,
		float32(boardOffsetX),
		float32(boardOffsetY),
		screenWidth,
		screenWidth,
		color.RGBA{70, 60, 94, 255},
		false,
	)

	gridColor := color.RGBA{100, 90, 130, 255}
	for i := 0; i <= boardSize; i++ {
		pos := float32(boardOffsetX + i*blockSize)
		// vertical
		vector.StrokeLine(screen, pos, float32(boardOffsetY), pos, float32(boardOffsetY+screenWidth), 1, gridColor, false)
		// horizontal
		hPos := float32(boardOffsetY + i*blockSize)
		vector.StrokeLine(screen, float32(boardOffsetX), hPos, float32(boardOffsetX+screenWidth), hPos, 1, gridColor, false)
	}

	// for _, point := range g.maze.carved {
	for y := range g.maze.height {
		for x := range g.maze.width {

			if !g.maze.grid[y][x] {
				continue
			}

			// path
			vector.FillRect(
				screen,
				float32(boardOffsetX+x*blockSize),
				float32(boardOffsetY+y*blockSize),
				blockSize,
				blockSize,
				color.White,
				false,
			)
		}
	}

	vector.FillRect(
		screen,
		float32(boardOffsetX+g.position.x*blockSize),
		float32(boardOffsetY+g.position.y*blockSize),
		blockSize,
		blockSize,
		color.RGBA{255, 0, 0, 255},
		false,
	)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Maze Visualizer")

	initialPosition := Point{1, 1}

	maze := NewMaze(initialPosition)

	g := &Game{
		position: initialPosition,
		maze:     maze,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
