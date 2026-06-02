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

func (g *Game) ChainDraw(screen *ebiten.Image) {

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.sim_Num))
	// fmt.Println(mgl64.Vec2{2, 2})

	g.chainPoints[0] = g.MousePos
	for i := 1; i < 30; i++ {
		g.chainPoints[i] = Constraint_dist(g.chainPoints[i], g.chainPoints[i-1], 20)
		vector.StrokeLine(screen, float32(g.chainPoints[i-1][0]), float32(g.chainPoints[i-1][1]), float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), 3, color.Black, true)
	}
	// chainstrokeopts := vector.StrokeOptions{Width: 2}

	// chainDrawPathOpts := vector.DrawPathOptions{AntiAlias: true}

	vector.FillCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 10, color.RGBA{255, 255, 255, 255}, true)

	// vector.StrokeLine(screen, float32(g.chainPoints[0][0]), float32(g.chainPoints[0][1]), float32(g.chainPoints[1][0]), float32(g.chainPoints[1][1]), 3, color.Black, true)

	// vector.StrokePath(screen, path, &chainstrokeopts, &chainDrawPathOpts)
	// vector.StrokeCircle(screen, float32(g.MousePos[0]), float32(g.MousePos[1]), 50, 1, color.RGBA{245, 40, 145, 100}, true)
	// vector.FillCircle(screen, float32(g.Point[0]), float32(g.Point[1]), 5, color.RGBA{245, 100, 145, 100}, true)

}
