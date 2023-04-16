//go:generate fyne bundle -o bundled.go directions.svg
//go:generate fyne bundle -a -o bundled.go temperature.svg

package main

import (
	"crypto/ed25519"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/fynelabs/fyneselfupdate"
	"github.com/fynelabs/selfupdate"
)

// selfManage turns on automatic updates
func selfManage(a fyne.App, w fyne.Window) {
	publicKey := ed25519.PublicKey{174, 168, 8, 15, 255, 74, 75, 30, 141, 156, 193, 40, 49, 94, 201, 184, 219, 38, 229, 172, 110, 43, 136, 88, 220, 230, 218, 185, 28, 155, 102, 230}

	httpSource := selfupdate.NewHTTPSource(nil, "https://geoffrey-artefacts.fynelabs.com/self-update/97/97e1436a-ed9b-451e-bc1d-0788875a2840/{{.OS}}-{{.Arch}}/{{.Executable}}{{.Ext}}")

	config := fyneselfupdate.NewConfig(a, w, httpSource, selfupdate.Schedule{FetchOnStart: true, Interval: time.Second * 15}, publicKey)

	_, err := selfupdate.Manage(config)
	if err != nil {
		log.Println("Error while setting up update manager: ", err)
		return
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Dash")

	temp := canvas.NewImageFromResource(theme.NewThemedResource(resourceTemperatureSvg))
	temp.FillMode = canvas.ImageFillContain
	temp.SetMinSize(fyne.NewSize(150, 80))
	dir := canvas.NewImageFromResource(theme.NewThemedResource(resourceDirectionsSvg))
	dir.FillMode = canvas.ImageFillContain
	dir.SetMinSize(fyne.NewSize(150, 80))

	w.SetContent(container.NewBorder(nil, nil,
		container.NewCenter(temp),
		container.NewCenter(dir),
		speedo(),
	))

	selfManage(a, w)
	w.ShowAndRun()
}
