package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"procedural-maze-visualizer/maze"
)

const (
	moleSpeed    = time.Second / 30
	screenWidth  = 720 / 3
	screenHeight = 1280 / 3
	blockSize    = 7
	boardSize    = screenWidth / blockSize
	boardOffsetX = 0
	boardOffsetY = (screenHeight - screenWidth) / 2
)

type Game struct {
	lastUpdate time.Time
	maze       *maze.Maze
}

func (g *Game) Update() error {
	if time.Since(g.lastUpdate) < moleSpeed {
		return nil
	}

	g.maze.Carve()

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
	for i := range boardSize {
		pos := float32(boardOffsetX + i*blockSize)
		// vertical
		vector.StrokeLine(screen, pos, float32(boardOffsetY), pos, float32(boardOffsetY+screenWidth), 1, gridColor, false)
		// horizontal
		hPos := float32(boardOffsetY + i*blockSize)
		vector.StrokeLine(screen, float32(boardOffsetX), hPos, float32(boardOffsetX+screenWidth), hPos, 1, gridColor, false)
	}

	for y := range g.maze.Height {
		for x := range g.maze.Width {

			if !g.maze.Grid[y][x] {
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
		float32(boardOffsetX+g.maze.Position.X*blockSize),
		float32(boardOffsetY+g.maze.Position.Y*blockSize),
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

	initialPosition := maze.Point{X: 1, Y: 1}
	// initialPosition := maze.Point{X: 1, Y: screenWidth / blockSize / 2}

	m := maze.NewMaze(initialPosition, screenWidth, blockSize)

	g := &Game{
		maze: m,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
