package main

import "encoding/json"

type MotionConfig struct {
	Mode          string  `json:"mode"`
	Segments      int     `json:"segments"`
	SegmentLength float64 `json:"segmentLength"`
	BodyWidth     float64 `json:"bodyWidth"`
	Amplitude     float64 `json:"amplitude"`
	Frequency     float64 `json:"frequency"`
	Speed         float64 `json:"speed"`
	Paused        bool    `json:"paused"`
}

func defaultMotionConfig() MotionConfig {
	return MotionConfig{Mode: "fish", Segments: 12, SegmentLength: 28, BodyWidth: 1, Amplitude: 18, Frequency: 1.2, Speed: 260}
}

func decodeMotionConfig(raw string) (MotionConfig, error) {
	c := defaultMotionConfig()
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return c, err
	}
	if c.Mode != "constraint" && c.Mode != "chain" && c.Mode != "fish" && c.Mode != "snake" && c.Mode != "lizard" {
		c.Mode = "fish"
	}
	if c.Segments < 6 {
		c.Segments = 6
	}
	if c.Segments > 24 {
		c.Segments = 24
	}
	if c.SegmentLength < 12 {
		c.SegmentLength = 12
	}
	if c.SegmentLength > 40 {
		c.SegmentLength = 40
	}
	if c.BodyWidth < .3 {
		c.BodyWidth = .3
	}
	if c.BodyWidth > 2 {
		c.BodyWidth = 2
	}
	if c.Amplitude < 0 {
		c.Amplitude = 0
	}
	if c.Amplitude > 50 {
		c.Amplitude = 50
	}
	if c.Frequency < .1 {
		c.Frequency = .1
	}
	if c.Frequency > 4 {
		c.Frequency = 4
	}
	if c.Speed < 40 {
		c.Speed = 40
	}
	if c.Speed > 720 {
		c.Speed = 720
	}
	return c, nil
}
