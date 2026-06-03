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

func Constraint_dist(point mgl64.Vec2, anchor mgl64.Vec2, distance float64) mgl64.Vec2 {
	delta := point.Sub(anchor)
	if delta.Len() < 0.0001 {
		return anchor
	}
	constraint_Dist := delta.Normalize()
	offset_Point := constraint_Dist.Mul(distance)

	return offset_Point.Add(anchor)
}

func Constraintupdate(anchor mgl64.Vec2, point *mgl64.Vec2) {
	nextpos := anchor.Sub(*point)
	if nextpos.Len() > Distance {
		*point = Constraint_dist(*point, anchor, Distance)
	}

}

func (g *Game) ConstraintDraw(screen *ebiten.Image) {

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.sim_Num))
	vector.StrokeCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 50, 1, color.RGBA{245, 40, 145, 100}, true)
	vector.FillCircle(screen, float32(g.Point[0]), float32(g.Point[1]), 5, color.RGBA{245, 100, 145, 100}, true)

}
