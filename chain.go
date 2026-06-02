package main

import (
	// "fmt"
	// "go/constant"
	// "fmt"
	"image/color"
	// "path"
	"strconv"

	// "github.com/go-gl/mathgl/mgl64"

	// "syscall/js"

	// "github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	Segments      = 30
	SegmentLength = 20
)

func (g *Game) ChainUpdate() {
	g.chainPoints[0] = g.MousePos
	for i := 1; i < Segments; i++ {
		g.chainPoints[i] = Constraint_dist(g.chainPoints[i], g.chainPoints[i-1], SegmentLength)
	}
}

func (g *Game) ChainDraw(screen *ebiten.Image) {

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.sim_Num))
	// fmt.Println(mgl64.Vec2{2, 2})

	for i := 1; i < Segments; i++ {
		vector.StrokeLine(screen, float32(g.chainPoints[i-1][0]), float32(g.chainPoints[i-1][1]), float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 3, color.Black, true)

		vector.FillCircle(screen, float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 3, color.RGBA{255, 255, 255, 255}, true)
	}

	vector.FillCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 3, color.RGBA{255, 255, 255, 255}, true)

}
