package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) ChainUpdate() {
	g.chainPoints[0] = g.MousePos
	for i := 1; i < g.Config.Segments; i++ {
		g.chainPoints[i] = Constraint_dist(g.chainPoints[i], g.chainPoints[i-1], g.Config.SegmentLength)
	}
}

func (g *Game) ChainDraw(screen *ebiten.Image) {

	for i := 1; i < g.Config.Segments; i++ {
		vector.StrokeLine(screen, float32(g.chainPoints[i-1][0]), float32(g.chainPoints[i-1][1]), float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 3, color.RGBA{82, 219, 204, 255}, true)
		vector.FillCircle(screen, float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 3, color.RGBA{245, 238, 218, 255}, true)
	}
	vector.FillCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 5, color.RGBA{246, 120, 95, 255}, true)

}
