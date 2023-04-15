//go:generate fyne bundle -o bundled.go directions.svg
//go:generate fyne bundle -a -o bundled.go temperature.svg

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	w := a.NewWindow("Dash")

	temp := canvas.NewImageFromResource(theme.NewThemedResource(resourceTemperatureSvg))
	temp.SetMinSize(fyne.NewSize(50, 50))
	dir := canvas.NewImageFromResource(theme.NewThemedResource(resourceDirectionsSvg))
	dir.SetMinSize(fyne.NewSize(50, 50))

	w.SetContent(container.NewGridWithColumns(3,
		container.NewCenter(temp),
		speedo(),
		container.NewCenter(dir)))
	w.ShowAndRun()
}
