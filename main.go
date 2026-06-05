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
	moleSpeed    = time.Second / 10
	screenWidth  = 720 / 3
	screenHeight = 1280 / 3
	mazeCols     = 31
	mazeRows     = 31
	paddingPx    = 4
	startDelay   = 3 * time.Second
)

var (
	blockSize   int
	mazeWidth   int
	mazeHeight  int
	mazeOffsetX int
	mazeOffsetY int
)

func init() {
	blockSize = min(
		(screenWidth-2*paddingPx)/mazeCols,
		(screenHeight-2*paddingPx)/mazeRows,
	)
	mazeWidth = mazeCols * blockSize
	mazeHeight = mazeRows * blockSize
	mazeOffsetX = (screenWidth - mazeWidth) / 2
	mazeOffsetY = (screenHeight - mazeHeight) / 2
}

type Game struct {
	startTime time.Time
	maze      *maze.Maze
	sounds    *Sounds
}

func (g *Game) Update() error {
	if time.Since(g.startTime) < startDelay {
		return nil
	}

	if time.Since(g.startTime)-startDelay < moleSpeed {
		return nil
	}

	if g.maze.Carve() {
		g.sounds.PlayCarve()
	} else if g.maze.Backtrack() {
		g.sounds.PlayBacktrack()
	}

	g.startTime = time.Now().Add(-startDelay)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 49, 59, 255})

	// bg
	vector.FillRect(
		screen,
		float32(mazeOffsetX),
		float32(mazeOffsetY),
		float32(mazeWidth),
		float32(mazeHeight),
		color.RGBA{70, 60, 94, 255},
		false,
	)

	gridColor := color.RGBA{100, 90, 130, 255}
	for i := range mazeCols {
		pos := float32(mazeOffsetX + i*blockSize)
		vector.StrokeLine(screen, pos, float32(mazeOffsetY), pos, float32(mazeOffsetY+mazeHeight), 1, gridColor, false)
	}
	for i := range mazeRows {
		hPos := float32(mazeOffsetY + i*blockSize)
		vector.StrokeLine(screen, float32(mazeOffsetX), hPos, float32(mazeOffsetX+mazeWidth), hPos, 1, gridColor, false)
	}

	for y := range g.maze.Height {
		for x := range g.maze.Width {

			if !g.maze.Grid[y][x] {
				continue
			}

			// path
			vector.FillRect(
				screen,
				float32(mazeOffsetX+x*blockSize),
				float32(mazeOffsetY+y*blockSize),
				float32(blockSize),
				float32(blockSize),
				color.White,
				false,
			)
		}
	}

	vector.FillRect(
		screen,
		float32(mazeOffsetX+g.maze.Position.X*blockSize),
		float32(mazeOffsetY+g.maze.Position.Y*blockSize),
		float32(blockSize),
		float32(blockSize),
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

	m := maze.NewMaze(initialPosition, mazeCols, mazeRows)

	g := &Game{
		startTime: time.Now(),
		maze:      m,
		sounds:    NewSounds(),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
