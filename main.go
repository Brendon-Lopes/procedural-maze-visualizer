package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
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

type Game struct {
	initialPosition Point
}

func (g *Game) Update() error {
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

	// debug grid
	gridColor := color.RGBA{100, 90, 130, 255}
	for i := 0; i <= boardSize; i++ {
		pos := float32(boardOffsetX + i*blockSize)
		// vertical
		vector.StrokeLine(screen, pos, float32(boardOffsetY), pos, float32(boardOffsetY+screenWidth), 1, gridColor, false)
		// horizontal
		hPos := float32(boardOffsetY + i*blockSize)
		vector.StrokeLine(screen, float32(boardOffsetX), hPos, float32(boardOffsetX+screenWidth), hPos, 1, gridColor, false)
	}

	initialX := boardOffsetX + g.initialPosition.x*blockSize
	initialY := boardOffsetY + g.initialPosition.y*blockSize

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

	g := &Game{
		initialPosition: Point{1, (boardSize - 1) / 2},
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
