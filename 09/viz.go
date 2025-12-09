package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/fogleman/gg"
)

type Visualizer struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	width  int
	height int
	// Coordinate mapping fields
	minX, minY    float64
	scale         float64
	offsetX       float64
	offsetY       float64
	polygonEdges  *[]Edge // Cache edges to draw every frame
	polygonPoints *[]Tile // Cache points to draw every frame
}

func NewVisualizer(filename string, edges *[]Edge, points *[]Tile) *Visualizer {
	// 1. Calculate Bounds of the polygon to fit screen
	minX, maxX, minY, maxY := 1e9, -1e9, 1e9, -1e9
	for _, p := range *points {
		if float64(p.x) < minX { minX = float64(p.x) }
		if float64(p.x) > maxX { maxX = float64(p.x) }
		if float64(p.y) < minY { minY = float64(p.y) }
		if float64(p.y) > maxY { maxY = float64(p.y) }
	}

	width, height := 800, 800
	
	// Calculate scale to fit with margin
	worldW := maxX - minX
	worldH := maxY - minY
	scaleX := float64(width) / worldW
	scaleY := float64(height) / worldH
	scale := scaleX
	if scaleY < scale { scale = scaleY }
	scale *= 0.9 // 10% margin

	// Center the drawing
	offX := (float64(width) - (worldW * scale)) / 2
	offY := (float64(height) - (worldH * scale)) / 2

	// 2. Start FFmpeg
	cmd := exec.Command("ffmpeg",
		"-y", "-f", "image2pipe", "-vcodec", "png", "-r", "60",
		"-i", "-", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-preset", "ultrafast",
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
		cmd:           cmd,
		stdin:         stdin,
		width:         width,
		height:        height,
		minX:          minX,
		minY:          minY,
		scale:         scale,
		offsetX:       offX,
		offsetY:       offY,
		polygonEdges:  edges,
		polygonPoints: points,
	}
}

// Transform World (x,y) to Screen (x,y)
func (v *Visualizer) toScreen(wx, wy int) (float64, float64) {
	sx := v.offsetX + (float64(wx)-v.minX)*v.scale
	sy := v.offsetY + (float64(wy)-v.minY)*v.scale
	return sx, sy
}

func (v *Visualizer) AddFrame(candT1, candT2, bestT1, bestT2 Tile, isNewBest bool, currentBestArea int) {
	dc := gg.NewContext(v.width, v.height)

	dc.SetHexColor("#000000") //black background
	dc.Clear()

	dc.SetRGB(0, 1, 0) // Green edges
	dc.SetLineWidth(2)
	for _, e := range *v.polygonEdges {
		x1, y1 := v.toScreen(e.start.x, e.start.y)
		x2, y2 := v.toScreen(e.end.x, e.end.y)
		dc.DrawLine(x1, y1, x2, y2)
	}
	dc.Stroke()

	// 3. Draw Corners 
	dc.SetRGB(1, 0, 0)
	for _, p := range *v.polygonPoints {
		px, py := v.toScreen(p.x, p.y)
		dc.DrawCircle(px, py, 3)
	}
	dc.Fill()

	if currentBestArea > 0 { // Draw best area rectangle
		bx1, by1 := v.toScreen(bestT1.x, bestT1.y)
		bx2, by2 := v.toScreen(bestT2.x, bestT2.y)
		bw := bx2 - bx1
		bh := by2 - by1
		
		dc.SetRGBA(1, 0.5, 0, 1.0) 
		dc.SetLineWidth(3)
		dc.DrawRectangle(bx1, by1, bw, bh)
		dc.Stroke()
	}

	rx1, ry1 := v.toScreen(candT1.x, candT1.y)
	rx2, ry2 := v.toScreen(candT2.x, candT2.y)
	w := rx2 - rx1
	h := ry2 - ry1

	if isNewBest {
		// New Record Flash: Bright White/Yellow Fill
		dc.SetRGBA(1, 1, 1, 0.8) 
		dc.DrawRectangle(rx1, ry1, w, h)
		dc.Fill()
		
		dc.SetRGB(1, 1, 1)
		dc.DrawString(fmt.Sprintf("NEW BEST AREA: %d", currentBestArea), 20, 50)
	} else {
		// Scanning: Faint Transparent Yellow
		dc.SetRGBA(1, 1, 0, 0.1) 
		dc.DrawRectangle(rx1, ry1, w, h)
		dc.Fill()
	}

	// Output frame
	dc.EncodePNG(v.stdin)
}

func (v *Visualizer) Close() {
	v.stdin.Close()
	v.cmd.Wait()
	fmt.Println("Video rendering complete.")
}