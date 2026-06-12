package main

import (
	// "fmt"
	// "go/constant"
	// "image/color"
	"image/color"
	"log"
	// "strconv"

	"github.com/go-gl/mathgl/mgl64"

	// "syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// "github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 360
	Distance     = 50
)

type Game struct {
	sim_Num     int
	xx          int
	yy          int
	Point       mgl64.Vec2
	Anchor      mgl64.Vec2
	MousePos    mgl64.Vec2
	chainPoints [30]mgl64.Vec2
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.sim_Num--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.sim_Num++
	}
	g.Anchor[0], g.Anchor[1] = float64(ScreenWidth/2), float64(ScreenHeight/2)

	x, y := ebiten.CursorPosition()
	g.MousePos[0] = float64(x)
	g.MousePos[1] = float64(y)

	switch g.sim_Num {
	case 0:
		Constraintupdate(g.MousePos, &g.Point)
	case 1:
		Constraintupdate(g.MousePos, &g.Point)
		g.ChainUpdate()
	case 2:

		Constraintupdate(g.MousePos, &g.Point)
		g.FishUpdate()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	bgColor := color.RGBA{R: 135, G: 206, B: 235, A: 255}

	// Fill the entire screen with the color
	screen.Fill(bgColor)
	// fmt.Println(gg)
	// fmt.Println(constraint_Dist.Len())
	// ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'g', -1, 64))
	switch g.sim_Num {
	case 0:

		g.ConstraintDraw(screen)
	case 1:
		g.ChainDraw(screen)
	case 2:
		g.FishDraw(screen)

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("HEllo World")
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
