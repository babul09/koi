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
	headInitialized         bool
	chainInitialized        bool
	mouseInitialized        bool
	heading                 float64
	layoutSegments          int
	layoutSegmentLength     float64
	chainAngles             [30]float64
	lizardLegInitialized    [4]bool
	lizardShoulders         [4]mgl64.Vec2
	lizardElbows            [4]mgl64.Vec2
	lizardFeet              [4]mgl64.Vec2
	lizardFootTargets       [4]mgl64.Vec2
	Point, Anchor, MousePos mgl64.Vec2
	chainPoints             [30]mgl64.Vec2
}

func (g *Game) ResetSimulation() {
	g.tick = 0
	g.headInitialized = false
	g.chainInitialized = false
	g.heading = 0
	g.layoutSegments = 0
	g.layoutSegmentLength = 0
	g.chainAngles = [30]float64{}
	g.lizardLegInitialized = [4]bool{}
	g.lizardShoulders = [4]mgl64.Vec2{}
	g.lizardElbows = [4]mgl64.Vec2{}
	g.lizardFeet = [4]mgl64.Vec2{}
	g.lizardFootTargets = [4]mgl64.Vec2{}
	g.chainPoints = [30]mgl64.Vec2{}
	g.Point = mgl64.Vec2{}
}

func (g *Game) Update() error {
	if g.Config.Mode == "" {
		g.Config = defaultMotionConfig()
	}
	x, y := ebiten.CursorPosition()
	if x >= 0 && y >= 0 && x < ScreenWidth && y < ScreenHeight {
		g.MousePos[0], g.MousePos[1] = float64(x), float64(y)
		g.mouseInitialized = true
	}
	if !g.mouseInitialized {
		g.MousePos = mgl64.Vec2{ScreenWidth / 2, ScreenHeight / 2}
		g.mouseInitialized = true
	}
	if !g.Config.Paused {
		g.tick += 1.0 / 60
		switch g.Config.Mode {
		case "constraint":
			Constraintupdate(g.MousePos, &g.Point)
		case "chain":
			g.ChainUpdate()
		case "fish", "snake", "lizard":
			g.FishUpdate()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 15, 11, 255})
	for x := 0; x < ScreenWidth; x += 40 {
		vector.StrokeLine(screen, float32(x), 0, float32(x), ScreenHeight, 1, color.RGBA{33, 39, 30, 80}, true)
	}
	for y := 0; y < ScreenHeight; y += 40 {
		vector.StrokeLine(screen, 0, float32(y), ScreenWidth, float32(y), 1, color.RGBA{33, 39, 30, 80}, true)
	}
	switch g.Config.Mode {
	case "constraint":
		g.ConstraintDraw(screen)
	case "chain":
		g.ChainDraw(screen)
	case "snake":
		g.AnimalDraw(screen, "snake")
	case "lizard":
		g.AnimalDraw(screen, "lizard")
	default:
		g.AnimalDraw(screen, "fish")
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
