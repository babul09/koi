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
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	FishSegments      = 10
	FishSegmentLength = 30
)

var FishWidth = []float64{45.33, 54.00, 56.00, 55.33, 51.33, 42.67, 34.00, 25.33, 21.33, 12.67}

var leftBody = make([]mgl64.Vec2, FishSegments)
var rightBody = make([]mgl64.Vec2, FishSegments)

func (g *Game) FishUpdate() {
	g.chainPoints[0] = g.MousePos
	for i := 1; i < FishSegments; i++ {
		g.chainPoints[i] = Constraint_dist(g.chainPoints[i], g.chainPoints[i-1], FishSegmentLength)
	}
	for i := 0; i < FishSegments; i++ {
		var dir mgl64.Vec2

		if i == 0 {
			delta := g.chainPoints[0].Sub(g.MousePos)
			if delta.Len() > 0.0001 {
				dir = delta.Normalize()
			}
		} else {
			delta := g.chainPoints[i].Sub(g.chainPoints[i-1])
			if delta.Len() > 0.0001 {
				dir = delta.Normalize()
			}
		}

		normal := mgl64.Vec2{
			-dir[1],
			dir[0],
		}

		leftBody[i] = g.chainPoints[i].Add(
			normal.Mul(FishWidth[i]),
		)

		rightBody[i] = g.chainPoints[i].Sub(
			normal.Mul(FishWidth[i]),
		)
	}
}

func (g *Game) FishDraw(screen *ebiten.Image) {

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.sim_Num))

	// fmt.Println(fishObject.bodyWidth)
	for i := 0; i < FishSegments; i++ {
		// vector.StrokeLine(screen, float32(g.chainPoints[i-1][0]), float32(g.chainPoints[i-1][1]), float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 3, color.Black, true)

		vector.FillCircle(screen, float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), float32(FishWidth[i]), color.RGBA{255, 255, 255, 20}, true)
	}

	for i := 0; i < FishSegments; i++ {
		vector.FillCircle(
			screen,
			float32(leftBody[i][0]),
			float32(leftBody[i][1]),
			4,
			color.RGBA{255, 0, 0, 255},
			true,
		)

		vector.FillCircle(
			screen,
			float32(rightBody[i][0]),
			float32(rightBody[i][1]),
			4,
			color.RGBA{0, 255, 0, 255},
			true,
		)

	}

	for i := 1; i < FishSegments; i++ {
		vector.StrokeLine(
			screen,
			float32(leftBody[i-1][0]),
			float32(leftBody[i-1][1]),
			float32(leftBody[i][0]),
			float32(leftBody[i][1]),
			float32(2),
			color.Black,
			true,
		)
		vector.StrokeLine(
			screen,
			float32(rightBody[i-1][0]),
			float32(rightBody[i-1][1]),
			float32(rightBody[i][0]),
			float32(rightBody[i][1]),
			float32(2),
			color.Black,
			true,
		)

	}

	// vector.FillCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 3, color.RGBA{255, 255, 255, 255}, true)

}
