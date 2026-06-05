package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	moleSpeed    = time.Second / 2
	screenWidth  = 300
	screenHeight = 480
	blockSize    = 20
	boardSize    = screenWidth / blockSize
	boardOffsetX = 0
	boardOffsetY = (screenHeight - screenWidth) / 2
)

type Maze struct {
	width  int
	height int
	grid   [][]bool
	rng    *rand.Rand
}

func NewMaze() *Maze {
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
	}
}

func (m *Maze) Carve(g *Game) {
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

			m.grid[ny][nx] = true
			m.grid[intermediateY][intermediateX] = true

			break
		}

		// TODO:
		// if it gets here, no valid block was found (dead end)
		// needs to backtrack
	}
}

func (m *Maze) Print() {
	for _, row := range m.grid {
		for _, cell := range row {
			if cell {
				fmt.Print("    ")
			} else {
				fmt.Print("██  ")
			}
		}
		fmt.Println()
		fmt.Println()
	}
}

type Point struct {
	x, y int
}

type Game struct {
	position      Point
	lastUpdate    time.Time
	mazeAlgorithm *Maze
}

func (g *Game) Update() error {
	if time.Since(g.lastUpdate) < moleSpeed {
		return nil
	}

	g.mazeAlgorithm.Carve(g)

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

	initialX := boardOffsetX + g.position.x*blockSize
	initialY := boardOffsetY + g.position.y*blockSize

	// mole
	vector.FillRect(
		screen,
		float32(initialX),
		float32(initialY),
		blockSize,
		blockSize,
		color.White,
		false,
	)

	ebitenutil.DebugPrint(screen, "initial X: "+strconv.Itoa(initialX))
	ebitenutil.DebugPrintAt(screen, "initial Y: "+strconv.Itoa(initialY), 0, 16)
	ebitenutil.DebugPrintAt(screen, "board size: "+strconv.Itoa(boardSize), 0, 32)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Maze Visualizer")

	maze := NewMaze()
	maze.Print()

	g := &Game{
		// initialPosition: Point{1, (boardSize - 1) / 2},
		position:      Point{1, 1},
		mazeAlgorithm: maze,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
