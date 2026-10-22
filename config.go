package main

import "encoding/json"

type MotionConfig struct {
	Mode          string  `json:"mode"`
	Segments      int     `json:"segments"`
	SegmentLength float64 `json:"segmentLength"`
	BodyWidth     float64 `json:"bodyWidth"`
	Amplitude     float64 `json:"amplitude"`
	Frequency     float64 `json:"frequency"`
	Paused        bool    `json:"paused"`
}

func defaultMotionConfig() MotionConfig {
	return MotionConfig{Mode: "fish", Segments: 10, SegmentLength: 24, BodyWidth: 1, Amplitude: 18, Frequency: 1.2}
}

func decodeMotionConfig(raw string) (MotionConfig, error) {
	c := defaultMotionConfig()
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return c, err
	}
	if c.Mode != "constraint" && c.Mode != "chain" && c.Mode != "fish" {
		c.Mode = "fish"
	}
	if c.Segments < 2 {
		c.Segments = 2
	}
	if c.Segments > 30 {
		c.Segments = 30
	}
	if c.SegmentLength < 5 {
		c.SegmentLength = 5
	}
	if c.SegmentLength > 60 {
		c.SegmentLength = 60
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
	return c, nil
}
