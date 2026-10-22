package main

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) FishUpdate() {
	segments := g.Config.Segments
	if !g.headInitialized || g.layoutSegments != segments || g.layoutSegmentLength != g.Config.SegmentLength {
		if !g.headInitialized {
			g.chainPoints[0] = mgl64.Vec2{ScreenWidth / 2, ScreenHeight / 2}
		}
		if !g.mouseInitialized {
			g.MousePos = g.chainPoints[0]
			g.mouseInitialized = true
		}
		g.heading = 0
		for i := 0; i < segments; i++ {
			g.chainAngles[i] = g.heading
			g.chainPoints[i] = g.chainPoints[0].Sub(mgl64.Vec2{math.Cos(g.heading), math.Sin(g.heading)}.Mul(float64(i) * g.Config.SegmentLength))
		}
		g.headInitialized = true
		g.layoutSegments = segments
		g.layoutSegmentLength = g.Config.SegmentLength
		g.lizardLegInitialized = [4]bool{}
	}
	head := g.chainPoints[0]
	targetDelta := g.MousePos.Sub(head)
	if targetDelta.Len() > 1 {
		step := math.Min(targetDelta.Len(), math.Min(16, g.Config.Speed/60))
		head = head.Add(targetDelta.Normalize().Mul(step))
		g.heading = math.Atan2(targetDelta[1], targetDelta[0])
	}
	g.chainPoints[0] = head
	g.chainAngles[0] = g.heading
	maxBend := math.Pi / 8
	if g.Config.Mode == "snake" {
		maxBend = math.Pi / 6
	} else if g.Config.Mode == "lizard" {
		maxBend = math.Pi / 10
	}
	for i := 1; i < segments; i++ {
		previous := g.chainPoints[i-1]
		current := g.chainPoints[i]
		curAngle := math.Atan2(previous[1]-current[1], previous[0]-current[0])
		previousAngle := g.chainAngles[i-1]
		wave := math.Sin(g.tick*g.Config.Frequency*2*math.Pi-float64(i)*.42) * (g.Config.Amplitude / 50) * maxBend * .45
		turn := wrapAngle(curAngle-previousAngle) + wave
		turn = math.Max(-maxBend, math.Min(maxBend, turn))
		angle := previousAngle + turn
		g.chainAngles[i] = angle
		g.chainPoints[i] = previous.Sub(mgl64.Vec2{math.Cos(angle), math.Sin(angle)}.Mul(g.Config.SegmentLength))
	}
	if g.Config.Mode == "lizard" {
		g.updateLizardLegs()
	}
}

func wrapAngle(angle float64) float64 {
	return math.Atan2(math.Sin(angle), math.Cos(angle))
}

func (g *Game) updateLizardLegs() {
	segments := g.Config.Segments
	indices := [2]int{segments / 3, segments * 2 / 3}
	stepLength := g.Config.SegmentLength
	for leg := 0; leg < 4; leg++ {
		pair := leg / 2
		index := indices[pair]
		index = max(1, min(segments-2, index))
		side := 1.0
		if leg%2 == 1 {
			side = -1
		}
		forward := g.chainPoints[index-1].Sub(g.chainPoints[index+1])
		if forward.Len() < .001 {
			forward = mgl64.Vec2{math.Cos(g.heading), math.Sin(g.heading)}
		} else {
			forward = forward.Normalize()
		}
		normal := mgl64.Vec2{-forward[1], forward[0]}
		shoulder := g.chainPoints[index].Add(normal.Mul(side * stepLength * .72))
		phaseOffset := 0.0
		if pair == 1 {
			phaseOffset = math.Pi
		}
		if side < 0 {
			phaseOffset += math.Pi
		}
		phase := g.tick*g.Config.Frequency*2*math.Pi + phaseOffset
		desired := shoulder.Add(normal.Mul(side * stepLength * .95)).Add(forward.Mul(math.Sin(phase) * stepLength * .32))
		if !g.lizardLegInitialized[leg] {
			g.lizardFootTargets[leg] = desired
			g.lizardFeet[leg] = desired
			g.lizardLegInitialized[leg] = true
		} else if desired.Sub(g.lizardFootTargets[leg]).Len() > stepLength*1.1 {
			g.lizardFootTargets[leg] = desired
		}
		g.lizardFeet[leg] = g.lizardFeet[leg].Add(g.lizardFootTargets[leg].Sub(g.lizardFeet[leg]).Mul(.28))
		g.lizardShoulders[leg] = shoulder
		g.lizardElbows[leg] = solveElbow(shoulder, g.lizardFeet[leg], stepLength*.92, side*float64(1-2*pair))
	}
}

func solveElbow(shoulder, foot mgl64.Vec2, linkLength, bendSide float64) mgl64.Vec2 {
	delta := foot.Sub(shoulder)
	distance := delta.Len()
	if distance < .001 {
		return shoulder.Add(mgl64.Vec2{0, bendSide * linkLength})
	}
	distance = math.Min(distance, linkLength*2-.001)
	direction := delta.Normalize()
	base := distance / 2
	height := math.Sqrt(math.Max(0, linkLength*linkLength-base*base))
	perpendicular := mgl64.Vec2{-direction[1], direction[0]}
	return shoulder.Add(direction.Mul(base)).Add(perpendicular.Mul(bendSide * height))
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
	fishProfile := []float64{68, 81, 84, 83, 77, 64, 51, 38, 32, 19}
	lizardProfile := []float64{52, 58, 40, 60, 68, 71, 65, 50, 28, 15, 11, 9, 7, 7}
	for i := 0; i < n; i++ {
		p := float64(i) / float64(n-1)
		var referenceWidth float64
		switch kind {
		case "snake":
			if i == 1 {
				referenceWidth = 80
			} else {
				referenceWidth = math.Max(17, 76-59*p)
			}
		case "lizard":
			referenceWidth = sampleProfile(lizardProfile, p)
		default:
			referenceWidth = sampleProfile(fishProfile, p)
		}
		widths[i] = referenceWidth * (g.Config.SegmentLength / 64) * g.Config.BodyWidth
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

func sampleProfile(profile []float64, t float64) float64 {
	position := t * float64(len(profile)-1)
	index := int(position)
	if index >= len(profile)-1 {
		return profile[len(profile)-1]
	}
	fraction := position - float64(index)
	return profile[index]*(1-fraction) + profile[index+1]*fraction
}

func smoothClosedPolygon(points []mgl64.Vec2, subdivisions int) []mgl64.Vec2 {
	if len(points) < 4 || subdivisions < 2 {
		return points
	}
	smoothed := make([]mgl64.Vec2, 0, len(points)*subdivisions)
	for i := range points {
		p0 := points[(i-1+len(points))%len(points)]
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		p3 := points[(i+2)%len(points)]
		for step := 0; step < subdivisions; step++ {
			t := float64(step) / float64(subdivisions)
			t2, t3 := t*t, t*t*t
			point := p1.Mul(2).
				Add(p2.Sub(p0).Mul(t)).
				Add(p0.Mul(2).Sub(p1.Mul(5)).Add(p2.Mul(4)).Sub(p3).Mul(t2)).
				Add(p0.Mul(-1).Add(p1.Mul(3)).Sub(p2.Mul(3)).Add(p3).Mul(t3)).
				Mul(.5)
			smoothed = append(smoothed, point)
		}
	}
	return smoothed
}

func (g *Game) AnimalDraw(screen *ebiten.Image, kind string) {
	n := g.Config.Segments
	outline, left, right := g.bodyGeometry(kind)
	var body, fin, ink color.RGBA
	switch kind {
	case "snake":
		body = color.RGBA{172, 57, 49, 255}
		fin = color.RGBA{216, 108, 83, 255}
		ink = color.RGBA{76, 33, 31, 255}
	case "lizard":
		body = color.RGBA{82, 121, 111, 255}
		fin = color.RGBA{123, 165, 143, 255}
		ink = color.RGBA{225, 225, 199, 255}
	default:
		body = color.RGBA{58, 124, 165, 255}
		fin = color.RGBA{129, 195, 215, 255}
		ink = color.RGBA{23, 56, 72, 255}
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
		for leg := range g.lizardFeet {
			shoulder, elbow, foot := g.lizardShoulders[leg], g.lizardElbows[leg], g.lizardFeet[leg]
			vector.StrokeLine(screen, float32(shoulder[0]), float32(shoulder[1]), float32(elbow[0]), float32(elbow[1]), float32(g.Config.SegmentLength*.17), ink, true)
			vector.StrokeLine(screen, float32(elbow[0]), float32(elbow[1]), float32(foot[0]), float32(foot[1]), float32(g.Config.SegmentLength*.12), body, true)
			vector.FillCircle(screen, float32(foot[0]), float32(foot[1]), 3, ink, true)
		}
	}
	drawPolygon(screen, smoothClosedPolygon(outline, 4), body, ink, 2)
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
}
