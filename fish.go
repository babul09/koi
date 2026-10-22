package main

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) FishUpdate() {
	if !g.headInitialized {
		g.chainPoints[0] = mgl64.Vec2{ScreenWidth / 2, ScreenHeight / 2}
		g.MousePos = g.chainPoints[0]
		g.headInitialized = true
	}
	head := g.chainPoints[0]
	targetDelta := g.MousePos.Sub(head)
	if targetDelta.Len() > 1 {
		targetHeading := math.Atan2(targetDelta[1], targetDelta[0])
		turn := math.Atan2(math.Sin(targetHeading-g.heading), math.Cos(targetHeading-g.heading))
		maxTurn := 4.0 / 60.0
		if turn > maxTurn {
			turn = maxTurn
		}
		if turn < -maxTurn {
			turn = -maxTurn
		}
		g.heading += turn
		step := g.Config.Speed / 60
		if step > targetDelta.Len() {
			step = targetDelta.Len()
		}
		head = head.Add(targetDelta.Normalize().Mul(step))
	}
	g.chainPoints[0] = head
	segments := g.Config.Segments
	spineAngle := g.heading + math.Pi
	amp := g.Config.Amplitude / 50
	maxBend := math.Pi / 8
	if g.Config.Mode == "lizard" {
		amp *= .24
	} else if g.Config.Mode == "snake" {
		maxBend = math.Pi / 5
	}
	for i := 1; i < segments; i++ {
		waveTurn := math.Sin(g.tick*g.Config.Frequency*2*math.Pi-float64(i)*.58) * amp * .46
		if waveTurn > maxBend {
			waveTurn = maxBend
		}
		if waveTurn < -maxBend {
			waveTurn = -maxBend
		}
		spineAngle += waveTurn
		g.chainPoints[i] = g.chainPoints[i-1].Add(mgl64.Vec2{math.Cos(spineAngle), math.Sin(spineAngle)}.Mul(g.Config.SegmentLength))
	}
}

func drawPolygon(screen *ebiten.Image, points []mgl64.Vec2, fill color.Color, stroke color.Color, strokeWidth float32) {
	if len(points) < 3 {
		return
	}
	var path vector.Path
	path.MoveTo(float32(points[0][0]), float32(points[0][1]))
	for _, p := range points[1:] {
		path.LineTo(float32(p[0]), float32(p[1]))
	}
	path.Close()
	var scale ebiten.ColorScale
	scale.ScaleWithColor(fill)
	vector.FillPath(screen, &path, nil, &vector.DrawPathOptions{AntiAlias: true, ColorScale: scale})
	if strokeWidth > 0 {
		for i := range points {
			j := (i + 1) % len(points)
			vector.StrokeLine(screen, float32(points[i][0]), float32(points[i][1]), float32(points[j][0]), float32(points[j][1]), strokeWidth, stroke, true)
		}
	}
}

func (g *Game) bodyGeometry(kind string) ([]mgl64.Vec2, []mgl64.Vec2, []mgl64.Vec2) {
	n := g.Config.Segments
	center := g.chainPoints[:n]
	left, right := make([]mgl64.Vec2, n), make([]mgl64.Vec2, n)
	widths := make([]float64, n)
	for i := 0; i < n; i++ {
		p := float64(i) / float64(n-1)
		var profile float64
		switch kind {
		case "snake":
			profile = .23 - .13*p
		case "lizard":
			profile = .31 + .21*math.Sin(math.Pi*p) - .20*p
		default:
			profile = .22 + .35*math.Sin(math.Pi*p) - .08*p
		}
		widths[i] = g.Config.SegmentLength * profile * g.Config.BodyWidth
		var tangent mgl64.Vec2
		if i == 0 {
			tangent = center[1].Sub(center[0])
		} else if i == n-1 {
			tangent = center[i].Sub(center[i-1])
		} else {
			tangent = center[i+1].Sub(center[i-1])
		}
		if tangent.Len() < .001 {
			tangent = mgl64.Vec2{math.Cos(g.heading + math.Pi), math.Sin(g.heading + math.Pi)}
		} else {
			tangent = tangent.Normalize()
		}
		normal := mgl64.Vec2{-tangent[1], tangent[0]}
		left[i] = center[i].Add(normal.Mul(widths[i]))
		right[i] = center[i].Sub(normal.Mul(widths[i]))
	}
	headDir := center[0].Sub(center[1])
	if headDir.Len() > .001 {
		headDir = headDir.Normalize()
	} else {
		headDir = mgl64.Vec2{math.Cos(g.heading), math.Sin(g.heading)}
	}
	front := center[0].Add(headDir.Mul(widths[0] * .7))
	outline := make([]mgl64.Vec2, 0, n*2+2)
	outline = append(outline, front)
	for i := 0; i < n; i++ {
		outline = append(outline, left[i])
	}
	for i := n - 1; i >= 0; i-- {
		outline = append(outline, right[i])
	}
	return outline, left, right
}

func (g *Game) AnimalDraw(screen *ebiten.Image, kind string) {
	n := g.Config.Segments
	outline, left, right := g.bodyGeometry(kind)
	var body, fin, ink color.RGBA
	switch kind {
	case "snake":
		body = color.RGBA{102, 164, 117, 255}
		fin = color.RGBA{164, 205, 133, 255}
		ink = color.RGBA{27, 67, 56, 255}
	case "lizard":
		body = color.RGBA{122, 174, 111, 255}
		fin = color.RGBA{187, 202, 126, 255}
		ink = color.RGBA{37, 72, 54, 255}
	default:
		body = color.RGBA{90, 126, 65, 255}
		fin = color.RGBA{166, 205, 98, 255}
		ink = color.RGBA{36, 61, 32, 255}
	}
	// Fins and limbs sit behind the torso so the spine reads as a single animal.
	if kind == "fish" {
		tailDir := g.chainPoints[n-1].Sub(g.chainPoints[n-2])
		if tailDir.Len() > .001 {
			tailDir = tailDir.Normalize()
		}
		tailNormal := mgl64.Vec2{-tailDir[1], tailDir[0]}
		tail := g.chainPoints[n-1]
		tailTip := tail.Add(tailDir.Mul(g.Config.SegmentLength * 1.1))
		notch := tail.Add(tailDir.Mul(g.Config.SegmentLength * .2))
		span := g.Config.SegmentLength * .65 * g.Config.BodyWidth
		drawPolygon(screen, []mgl64.Vec2{g.chainPoints[n-2], tail, tailTip.Add(tailNormal.Mul(span)), notch, tailTip.Sub(tailNormal.Mul(span))}, fin, ink, 1.5)
		for _, idx := range []int{n / 3} {
			if idx > 0 && idx < n {
				p := g.chainPoints[idx]
				tangent := g.chainPoints[idx+1].Sub(g.chainPoints[idx-1])
				if tangent.Len() > .001 {
					tangent = tangent.Normalize()
				}
				normal := mgl64.Vec2{-tangent[1], tangent[0]}
				for _, side := range []float64{-1, 1} {
					drawPolygon(screen, []mgl64.Vec2{p.Add(normal.Mul(side * 5)), p.Add(normal.Mul(side * g.Config.SegmentLength * .8)).Add(tangent.Mul(g.Config.SegmentLength * .75)), p.Add(normal.Mul(side * 2)).Add(tangent.Mul(g.Config.SegmentLength * .3))}, fin, ink, 1)
				}
			}
		}
		// Dorsal fin.
		i := n * 2 / 5
		j := n * 3 / 5
		p1 := g.chainPoints[i]
		p2 := g.chainPoints[j]
		normal := left[i].Sub(g.chainPoints[i])
		if normal.Len() > .001 {
			normal = normal.Normalize()
		}
		drawPolygon(screen, []mgl64.Vec2{p1, p2, p2.Add(normal.Mul(g.Config.SegmentLength * .38)), p1.Add(normal.Mul(g.Config.SegmentLength * .18))}, fin, ink, 1)
	} else if kind == "lizard" {
		for _, idx := range []int{n / 3, 2 * n / 3} {
			if idx < 1 || idx >= n-1 {
				continue
			}
			tangent := g.chainPoints[idx+1].Sub(g.chainPoints[idx-1])
			if tangent.Len() > .001 {
				tangent = tangent.Normalize()
			}
			normal := mgl64.Vec2{-tangent[1], tangent[0]}
			phase := g.tick*g.Config.Frequency*2*math.Pi + float64(idx)
			for _, side := range []float64{-1, 1} {
				hip := g.chainPoints[idx].Add(normal.Mul(side * g.Config.SegmentLength * .18))
				elbow := hip.Add(normal.Mul(side * g.Config.SegmentLength * .53)).Add(tangent.Mul(math.Sin(phase) * g.Config.SegmentLength * .16))
				foot := elbow.Add(normal.Mul(side * g.Config.SegmentLength * .32)).Add(tangent.Mul(g.Config.SegmentLength * .15))
				vector.StrokeLine(screen, float32(hip[0]), float32(hip[1]), float32(elbow[0]), float32(elbow[1]), float32(g.Config.SegmentLength*.13), ink, true)
				vector.StrokeLine(screen, float32(elbow[0]), float32(elbow[1]), float32(foot[0]), float32(foot[1]), float32(g.Config.SegmentLength*.10), fin, true)
				vector.FillCircle(screen, float32(foot[0]), float32(foot[1]), 3, ink, true)
			}
		}
	}
	drawPolygon(screen, outline, body, ink, 2)
	// Belly stripe on fish; repeating dorsal marks on snakes.
	if kind == "fish" {
		for i := 2; i < n-2; i++ {
			vector.StrokeLine(screen, float32(left[i][0]), float32(left[i][1]), float32(right[i][0]), float32(right[i][1]), 1, color.RGBA{153, 218, 217, 95}, true)
		}
	}
	if kind == "snake" {
		for i := 2; i < n-1; i += 2 {
			vector.StrokeLine(screen, float32(left[i][0]), float32(left[i][1]), float32(right[i][0]), float32(right[i][1]), 2, fin, true)
		}
	}
	// Eyes sit just behind the blunt head tip and follow the spine normal.
	headAxis := g.chainPoints[0].Sub(g.chainPoints[1])
	if headAxis.Len() > .001 {
		headAxis = headAxis.Normalize()
	}
	normal := mgl64.Vec2{-headAxis[1], headAxis[0]}
	eyeCenter := g.chainPoints[0].Add(headAxis.Mul(g.Config.SegmentLength * .12))
	eyeOffset := math.Max(3, g.Config.SegmentLength*.16*g.Config.BodyWidth)
	for _, side := range []float64{-1, 1} {
		eye := eyeCenter.Add(normal.Mul(side * eyeOffset))
		vector.FillCircle(screen, float32(eye[0]), float32(eye[1]), float32(g.Config.SegmentLength*.105), color.RGBA{248, 243, 220, 255}, true)
		pupil := eye.Add(headAxis.Mul(2))
		vector.FillCircle(screen, float32(pupil[0]), float32(pupil[1]), float32(g.Config.SegmentLength*.052), color.RGBA{17, 30, 34, 255}, true)
	}
	if kind == "fish" {
		gill := g.chainPoints[2]
		vector.StrokeLine(screen, float32(gill[0]), float32(gill[1]), float32(gill[0]+normal[0]*5), float32(gill[1]+normal[1]*5), 1.5, ink, true)
	}
	// Spine markers keep the live constraint rig visible through the silhouette.
	spineColor := color.RGBA{199, 245, 107, 225}
	for i := 0; i < n; i++ {
		if i > 0 {
			start, end := g.chainPoints[i-1], g.chainPoints[i]
			steps := int(end.Sub(start).Len() / 8)
			for step := 0; step <= steps; step += 2 {
				t := float64(step) / float64(max(steps, 1))
				point := start.Add(end.Sub(start).Mul(t))
				vector.FillCircle(screen, float32(point[0]), float32(point[1]), 1.5, color.RGBA{199, 245, 107, 145}, true)
			}
		}
		radius := float32(2.4)
		if i == 0 || i == n-1 { radius = 3.5 }
		vector.FillCircle(screen, float32(g.chainPoints[i][0]), float32(g.chainPoints[i][1]), radius, spineColor, true)
	}
}
