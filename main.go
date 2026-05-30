package main

import (
	// "github.com/go-gl/mathgl/mgl64
	"image/color"
	"log"
	"strconv"
	// "syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 360
)

type Game struct {
	score int
	xx    int
	yy    int
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.score++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.xx, g.yy = ebiten.CursorPosition()

	// ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'g', -1, 64))
	ebitenutil.DebugPrint(screen, strconv.Itoa(g.score))
	vector.StrokeCircle(screen, float32(g.xx), float32(g.yy), 20, 1, color.RGBA{245, 40, 145, 100}, true)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("HEllo World")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
