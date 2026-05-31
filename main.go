package main

import (
	// "fmt"
	// "go/constant"
	"image/color"
	"log"
	"strconv"

	"github.com/go-gl/mathgl/mgl64"

	// "syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 360
	Distance     = 50
)

type Game struct {
	sim_Num int
	xx      int
	yy      int
	Point   mgl64.Vec2
	Anchor  mgl64.Vec2
}

type Points struct {
	rope mgl64.Vec2
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.sim_Num++
	}
	g.Anchor[0], g.Anchor[1] = float64(ScreenWidth/2), float64(ScreenHeight/2)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.xx, g.yy = ebiten.CursorPosition()
	act_Point := g.Chain(g.xx, g.yy)
	// fmt.Println(gg)
	// fmt.Println(constraint_Dist.Len())
	// ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'g', -1, 64))
	switch g.sim_Num {
	case 0:

		ebitenutil.DebugPrint(screen, strconv.Itoa(g.sim_Num))
		vector.StrokeCircle(screen, float32(g.xx), float32(g.yy), 20, 1, color.RGBA{245, 40, 145, 100}, true)
		vector.FillCircle(screen, float32(g.Anchor[0]), float32(g.Anchor[1]), 5, color.RGBA{245, 100, 145, 100}, true)
		vector.FillCircle(screen, float32(g.Anchor[0]+act_Point[0]), float32(g.Anchor[1]+act_Point[1]), 5, color.RGBA{245, 100, 145, 100}, true)
	}
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
