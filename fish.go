package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

var Front mgl64.Vec2

func (g *Game) FishUpdate() {
	n := g.Config.Segments
	g.chainPoints[0] = g.MousePos
	for i := 1; i < n; i++ {
		base := Constraint_dist(g.chainPoints[i], g.chainPoints[i-1], g.Config.SegmentLength)
		delta := base.Sub(g.chainPoints[i-1])
		phase := g.tick*g.Config.Frequency*2*math.Pi - float64(i)*.48
		base[1] += math.Sin(phase) * g.Config.Amplitude * float64(i) / float64(n)
		if delta.Len() > 0.001 {
			base[1] += math.Sin(phase) * g.Config.Amplitude * .16
		}
		g.chainPoints[i] = base
	}
	for i := 0; i < n; i++ {
		if i == 0 {
			Front = g.chainPoints[i]
		}
	}
}

func (g *Game) FishDraw(screen *ebiten.Image) {
	n := g.Config.Segments
	left := make([]mgl64.Vec2, n)
	right := make([]mgl64.Vec2, n)
	for i := 0; i < n; i++ {
		var dir mgl64.Vec2
		if i+1 < n {
			dir = g.chainPoints[i+1].Sub(g.chainPoints[i])
		} else if i > 0 {
			dir = g.chainPoints[i].Sub(g.chainPoints[i-1])
		}
		if dir.Len() > .001 {
			dir = dir.Normalize()
		}
		normal := mgl64.Vec2{-dir[1], dir[0]}
		width := (42*(1-float64(i)/float64(n)) + 5) * g.Config.BodyWidth
		left[i] = g.chainPoints[i].Add(normal.Mul(width))
		right[i] = g.chainPoints[i].Sub(normal.Mul(width))
	}
	for i := 1; i < n; i++ {
		vector.StrokeLine(screen, float32(left[i-1][0]), float32(left[i-1][1]), float32(left[i][0]), float32(left[i][1]), 2, color.RGBA{246, 120, 95, 255}, true)
		vector.StrokeLine(screen, float32(right[i-1][0]), float32(right[i-1][1]), float32(right[i][0]), float32(right[i][1]), 2, color.RGBA{246, 120, 95, 255}, true)
		vector.StrokeLine(screen, float32(left[i-1][0]), float32(left[i-1][1]), float32(right[i-1][0]), float32(right[i-1][1]), 1, color.RGBA{246, 120, 95, 140}, true)
	}
	vector.StrokeLine(screen, float32(left[n-1][0]), float32(left[n-1][1]), float32(right[n-1][0]), float32(right[n-1][1]), 1, color.RGBA{246, 120, 95, 140}, true)
	for i := 0; i < n; i++ {
		vector.FillCircle(screen, float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 2, color.RGBA{247, 198, 132, 210}, true)
	}
	vector.FillCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 4, color.RGBA{246, 120, 95, 255}, true)
}
