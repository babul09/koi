package main

import (
	"image/color"
	"log"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 600
	Distance     = 50
)

type Game struct {
	Config                  MotionConfig
	tick                    float64
	Point, Anchor, MousePos mgl64.Vec2
	chainPoints             [30]mgl64.Vec2
}

func (g *Game) Update() error {
	if g.Config.Mode == "" {
		g.Config = defaultMotionConfig()
	}
	x, y := ebiten.CursorPosition()
	g.MousePos[0], g.MousePos[1] = float64(x), float64(y)
	if !g.Config.Paused {
		g.tick += 1.0 / 60
	}
	switch g.Config.Mode {
	case "constraint":
		Constraintupdate(g.MousePos, &g.Point)
	case "chain":
		g.ChainUpdate()
	case "fish":
		g.FishUpdate()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{17, 29, 39, 255})
	for x := 0; x < ScreenWidth; x += 40 {
		vector.StrokeLine(screen, float32(x), 0, float32(x), ScreenHeight, 1, color.RGBA{42, 58, 66, 110}, true)
	}
	for y := 0; y < ScreenHeight; y += 40 {
		vector.StrokeLine(screen, 0, float32(y), ScreenWidth, float32(y), 1, color.RGBA{42, 58, 66, 110}, true)
	}
	switch g.Config.Mode {
	case "constraint":
		g.ConstraintDraw(screen)
	case "chain":
		g.ChainDraw(screen)
	default:
		g.FishDraw(screen)
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) { return ScreenWidth, ScreenHeight }

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Motion DNA")
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g := &Game{Config: defaultMotionConfig()}
	installBrowserBridge(g)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
