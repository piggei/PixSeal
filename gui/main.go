package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/pj/pixseal/internal/buildinfo"
)

func main() {
	a := app.NewWithID("io.github.piggei.pixseal")
	a.Settings().SetTheme(pixTheme{})
	w := a.NewWindow("PixSeal " + buildinfo.String())
	NewView(a, w, nativeEngine{})
	w.Resize(fyne.NewSize(390, 820))
	w.ShowAndRun()
}
