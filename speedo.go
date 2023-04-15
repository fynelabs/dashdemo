package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type dialLayout struct {
	needle *canvas.Line
	pips   [121]*canvas.Line
	face   *canvas.Circle
	cover  *canvas.Rectangle
	speed  *canvas.Text

	canvas fyne.CanvasObject
	stop   bool
	value float64
}

func (c *dialLayout) rotate(hand *canvas.Line, middle fyne.Position, facePosition float64, offset, length float32) {
	rotation := math.Pi * 1.5 / 120 * facePosition - math.Pi/4*3
	x2 := length * float32(math.Sin(rotation))
	y2 := -length * float32(math.Cos(rotation))

	offX := float32(0)
	offY := float32(0)
	if offset > 0 {
		offX += offset * float32(math.Sin(rotation))
		offY += -offset * float32(math.Cos(rotation))
	}

	hand.Position1 = fyne.NewPos(middle.X+offX, middle.Y+offY)
	hand.Position2 = fyne.NewPos(middle.X+offX+x2, middle.Y+offY+y2)
}

func (c *dialLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	c.setPosition(c.value, size)
}

func (c *dialLayout) setPosition(v float64, size fyne.Size) {
	c.value = v
	diameter := fyne.Min(size.Width, size.Height)
	radius := diameter / 2
	stroke := diameter / 50
	midStroke := diameter / 80
	smallStroke := diameter / 200

	size = fyne.NewSize(diameter, diameter)
	middle := fyne.NewPos(size.Width/2, size.Height/2)
	topleft := fyne.NewPos(middle.X-radius, middle.Y-radius)

	c.face.Move(topleft)
	c.face.Resize(size)
	c.cover.Move(fyne.NewPos(0, middle.Y+radius/7*5))
	c.cover.Resize(fyne.NewSize(size.Width, size.Height/6))
	c.speed.Move(topleft)
	c.speed.Resize(size)
	c.speed.TextSize = size.Height/3
	c.speed.Text = fmt.Sprintf("%d", int(v))

	c.needle.StrokeWidth = stroke
	c.rotate(c.needle, middle, v, radius*.2, radius*.75)
	c.face.StrokeWidth = smallStroke

	for i, p := range c.pips {
		if i % 10 == 0 {
			c.rotate(p, middle, float64(i), radius/4*3, radius/4)
			p.StrokeWidth = midStroke
		} else {
			c.rotate(p, middle, float64(i), radius/8*7, radius/8)
			p.StrokeWidth = smallStroke
		}
	}
}

func (c *dialLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(150, 150)
}

func (c *dialLayout) render() *fyne.Container {
	c.face = &canvas.Circle{StrokeColor: theme.DisabledColor(), StrokeWidth: 1}
	c.cover = &canvas.Rectangle{FillColor: theme.BackgroundColor()}
	c.needle = &canvas.Line{StrokeColor: theme.ErrorColor(), StrokeWidth: 7}
	c.speed = &canvas.Text{Text: "60", Color: theme.ForegroundColor(), TextSize: 52}
	c.speed.TextStyle.Monospace = true
	c.speed.Alignment = fyne.TextAlignCenter

	container := container.NewWithoutLayout(c.face, c.cover)
	for i, _ := range c.pips {
		pip := &canvas.Line{StrokeColor: theme.DisabledColor(), StrokeWidth: 1}
		if i == 0 {
			pip.StrokeColor = theme.ForegroundColor()
		} else if i >= 100 && i < 110 {
			pip.StrokeColor = theme.WarningColor()
		} else if i >= 110 {
			pip.StrokeColor = theme.ErrorColor()
		}
		container.Add(pip)
		c.pips[i] = pip
	}
	container.Objects = append(container.Objects, c.needle, c.speed)
	container.Layout = c

	c.canvas = container
	return container
}

func (c *dialLayout) animate(co fyne.CanvasObject) {
	tick := time.NewTicker(time.Second)
	go func() {
		for !c.stop {
			start := c.value
			stop := rand.Float64()*115
			fyne.NewAnimation(time.Second, func(v float32) {
				val := start + (stop-start)*float64(v)
				c.setPosition(val, co.Size())
				c.needle.Refresh()
				c.speed.Refresh()
			}).Start()
			<-tick.C
		}
	}()
}

func (c *dialLayout) applyTheme(_ fyne.Settings) {
	c.face.StrokeColor = theme.DisabledColor()
	c.needle.StrokeColor = theme.ErrorColor()
	c.speed.Color = theme.ForegroundColor()

	for _, p := range c.pips {
		p.StrokeColor = theme.DisabledColor()
	}
}

// speedo loads a speedo example window for the specified app context
func speedo() fyne.CanvasObject {
	s := &dialLayout{}

	content := s.render()
	go s.animate(content)

	listener := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			s.applyTheme(settings)
		}
	}()

	return content
}
