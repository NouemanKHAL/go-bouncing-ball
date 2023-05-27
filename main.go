package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	LINE_RETURN byte = '\n'
	SPACE       byte = ' '
	UNDERSCORE  byte = '_'
	CARET       byte = '^'
	C           byte = 'C'
	STAR        byte = '*'
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Pixel int

const (
	EMPTY Pixel = iota
	FILL
)

type Point struct {
	X, Y int
}

type Display struct {
	Width  int
	Height int
	grid   [][]Pixel
	writer io.Writer
}

func NewDisplay(height, width int, w io.Writer) *Display {
	display := make([][]Pixel, height)
	for i := 0; i < height; i++ {
		display[i] = make([]Pixel, width)
	}
	return &Display{
		Width:  width,
		Height: height,
		grid:   display,
		writer: w,
	}
}

func (d *Display) Fill(b Pixel) {
	for i := 0; i < d.Height; i++ {
		for j := 0; j < d.Width; j++ {
			d.grid[i][j] = b
		}
	}
}

func (d *Display) Back() {
	d.writer.Write([]byte(fmt.Sprintf("\x1b[%dD", d.Width)))
	d.writer.Write([]byte(fmt.Sprintf("\x1b[%dA", d.Height/2)))
}

func compressPixels(a, b Pixel) byte {
	tab := []byte{
		SPACE, UNDERSCORE, CARET, C,
	}
	return tab[int(a)*2+int(b)]
}

func (d *Display) Show() {
	row := make([]byte, d.Width)
	for i := 0; i < d.Height/2; i++ {
		for j := 0; j < d.Width; j++ {
			row[j] = compressPixels(d.grid[2*i][j], d.grid[2*i+1][j])
		}
		d.writer.Write(row)
		d.writer.Write([]byte{LINE_RETURN})
	}
}

func (d *Display) Circle(p Point, r int) {
	for i := p.X - r; i <= p.X+r; i++ {
		for j := p.Y - r; j <= p.Y+r; j++ {
			if i < 0 || i >= d.Height || j < 0 || j >= d.Width {
				continue
			}
			// dx and dy could be negative but it doesn't matter since we use their square
			dx := i - p.X
			dy := j - p.Y

			if dx*dx+dy*dy <= r*r {
				d.grid[i][j] = FILL
			}
		}
	}
}

func main() {
	height := 32
	width := 64
	radius := height / 4

	var velocity float32 = 0
	var gravity float32 = 250.0
	var dt float32 = 1.0 / 30

	display := NewDisplay(height, width, os.Stdout)
	start := Point{0, 0}
	for start.Y < display.Width+2*radius {
		velocity = velocity + gravity*dt
		start.X += int(velocity * dt)
		start.Y += 1
		if start.X > display.Height-radius {
			start.X = display.Height - radius
			velocity *= -0.68
		}

		display.Fill(EMPTY)
		display.Circle(start, radius)
		display.Show()
		display.Back()

		time.Sleep(time.Duration((1000 / 30) * time.Millisecond))
	}
}
