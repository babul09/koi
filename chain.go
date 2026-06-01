package main

import (
	// "fmt"
	// "go/constant"
	"image/color"
	"strconv"

	"github.com/go-gl/mathgl/mgl64"

	// "syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Constraint_dist(point mgl64.Vec2, anchor mgl64.Vec2, distance int) mgl64.Vec2 {
	// var P mgl64.Vec2
	// P[0], P[1] = float64(xx), float64(yy)
	delta := point.Sub(g.Anchor)
	constraint_Dist := delta.Normalize()
	offset_Point := constraint_Dist.Mul(Distance)

	return offset_Point
}

func (g *Game) Cmain(screen *ebiten.Image) {

	x, y := ebiten.CursorPosition()
	var anchor mgl64.Vec2
	anchor[0] = float64(x)
	anchor[1] = float64(y)
	point := mgl64.Vec2{0, 0}
	act_Point := g.Constraint_dist(point, anchor, Distance)

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.sim_Num))
	vector.StrokeCircle(screen, float32(x), float32(y), 50, 1, color.RGBA{245, 40, 145, 100}, true)
	vector.FillCircle(screen, float32(point[0]), float32(point[1]), 5, color.RGBA{245, 100, 145, 100}, true)
	vector.FillCircle(screen, float32(point[0]+act_Point[0]), float32(point[1]+act_Point[1]), 5, color.RGBA{245, 100, 145, 100}, true)

}
