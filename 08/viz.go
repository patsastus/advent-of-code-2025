package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os/exec"

	"github.com/fogleman/gg"
)
//----WARNING:nearly fully vibe-coded---
// --- 1. Data Structures ---

type VizPoint struct {
	X, Y, Z float64
	Size    int
}

type VizEdge struct {
	FromIndex int
	ToIndex   int
}

// --- 2. The Visualizer Struct ---

type Visualizer struct {
	cmd        *exec.Cmd
	stdin      io.WriteCloser
	width      int
	height     int
	scale      float64
	center     float64
	frameCount int
}

// NewVisualizer starts FFmpeg in the background
func NewVisualizer(filename string) *Visualizer {
	// Standard FFmpeg command for MP4 output
	cmd := exec.Command("ffmpeg",
		"-y", "-f", "image2pipe", "-vcodec", "png", "-r", "60",
		"-i", "-", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-preset", "fast",
		filename,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	return &Visualizer{
		cmd:    cmd,
		stdin:  stdin,
		width:  800,
		height: 800,
		scale:  600.0,   // Adjust zoom
		center: 50000.0, // Center of your 0-100k coordinate space
	}
}

func (v *Visualizer) AddFrame(points []VizPoint, edges []VizEdge, totalNodes int) {
	v.frameCount++

	dc := gg.NewContext(v.width, v.height)

	// 1. Background (Black)
	dc.SetHexColor("#000000")
	dc.Clear()

	// 2. Projection (Same as before)
	angle := float64(v.frameCount) * 0.0005
	scaleFactor := v.scale / (v.center * 2)

	projected := make([]struct {
		X, Y float64
		Size int
	}, len(points))
	for i, p := range points {
		px := p.X - v.center
		py := p.Y - v.center
		pz := p.Z - v.center

		rx := px*math.Cos(angle) - pz*math.Sin(angle)

		screenX := (float64(v.width) / 2) + (rx * scaleFactor)
		screenY := (float64(v.height) / 2) + (py * scaleFactor)

		projected[i] = struct {
			X, Y float64
			Size int
		}{screenX, screenY, p.Size}
	}

	// 3. Draw Edges
	// We make them faint gray so they don't distract from the color heat map
	dc.SetRGBA(1, 1, 1, 0.2)
	dc.SetLineWidth(3)
	for _, e := range edges {
		p1 := projected[e.FromIndex]
		p2 := projected[e.ToIndex]
		dc.DrawLine(p1.X, p1.Y, p2.X, p2.Y)
	}
	dc.Stroke()

	// 4. Draw Nodes (Color by Size)
	for _, p := range projected {
		r, g, b := colorFromSize(p.Size, totalNodes)
		dc.SetRGB(r, g, b)
		dc.DrawCircle(p.X, p.Y, 4)
		dc.Fill()
	}

	dc.EncodePNG(v.stdin)
}

// Logarithmic Gradient: Green -> Yellow -> Red
func colorFromSize(size, max int) (float64, float64, float64) {
	// 1. Safety check
	if size <= 1 {
		return 0, 0.6, 0
	} // Dark Green for unconnected nodes

	// 2. Logarithmic Normalization
	// We map the range [1 ... max] to [0.0 ... 1.0]
	// but on a Log scale.
	minLog := math.Log(1.0) // is 0
	maxLog := math.Log(float64(max))
	valLog := math.Log(float64(size))

	// t represents our progress from 0 to 1 based on Orders of Magnitude
	t := (valLog - minLog) / (maxLog - minLog)

	// 3. Color Interpolation (Green -> Yellow -> Red)
	r := 0.0
	g := 0.0

	// We shift the curve slightly so 't' isn't too aggressive immediately,
	// but much faster than linear.
	if t < 0.5 {
		// Green to Yellow
		r = t * 2.0
		g = 1.0
	} else {
		// Yellow to Red
		r = 1.0
		g = 1.0 - ((t - 0.5) * 2.0)
	}

	return r, g, 0.0
}

func colorFromID(id int) (float64, float64, float64) {
	if id == -1 {
		return 0.5, 0.5, 0.5
	} // Gray for unconnected
	// Golden ratio hash for distinct colors
	h := float64(id) * 0.618033988749895
	h -= math.Floor(h)
	return HSVToRGB(h, 0.8, 0.9) 
}

func HSVToRGB(h, s, v float64) (r, g, b float64) {
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	switch int(i) % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	}
	return 0, 0, 0
}

// Close finishes the video file
func (v *Visualizer) Close() {
	v.stdin.Close()
	v.cmd.Wait()
	fmt.Println("Video rendering complete.")
}
