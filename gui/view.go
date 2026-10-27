package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pj/pixseal/internal/appservice"
)

type Engine interface {
	Load(io.Reader, string) (appservice.Document, error)
	Embed(appservice.Document, string, string) (appservice.Result, error)
	Extract(appservice.Document, string, string, int, int) (appservice.Result, error)
}
type nativeEngine struct{}

func (nativeEngine) Load(r io.Reader, n string) (appservice.Document, error) {
	return appservice.Load(r, n)
}
func (nativeEngine) Embed(d appservice.Document, m, k string) (appservice.Result, error) {
	return appservice.Embed(d, m, k)
}
func (nativeEngine) Extract(d appservice.Document, k, m string, w, h int) (appservice.Result, error) {
	return appservice.Extract(d, k, m, w, h)
}

type View struct {
	app                                                            fyne.App
	window                                                         fyne.Window
	engine                                                         Engine
	doc                                                            appservice.Document
	last                                                           appservice.Result
	busy                                                           atomic.Bool
	generation                                                     atomic.Uint64
	key, message, width, height, output                            *widget.Entry
	choose, read, create, savePNG, saveReport, usePNG, copyMessage *widget.Button
	mode                                                           *widget.Select
	file, status, elapsed, byteCount                               *widget.Label
	preview                                                        *canvas.Image
	progress                                                       *widget.ProgressBarInfinite
	tabs                                                           *container.AppTabs
}

func wrap(s string) *widget.Label { l := widget.NewLabel(s); l.Wrapping = fyne.TextWrapWord; return l }
func field(label string, item fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), item)
}
func NewView(a fyne.App, w fyne.Window, e Engine) *View {
	v := &View{app: a, window: w, engine: e}
	v.key = widget.NewPasswordEntry()
	v.key.SetPlaceHolder("Almeno 8 byte · non viene salvata")
	v.message = widget.NewMultiLineEntry()
	v.message.SetMinRowsVisible(3)
	v.message.SetPlaceHolder("Il messaggio da nascondere…")
	v.byteCount = wrap("0 / 64 byte UTF-8")
	v.message.OnChanged = func(s string) { v.byteCount.SetText(fmt.Sprintf("%d / 64 byte UTF-8", len([]byte(s)))) }
	v.width = widget.NewEntry()
	v.height = widget.NewEntry()
	v.width.SetText(strconv.Itoa(a.Preferences().IntWithFallback("canonical_width", 1632)))
	v.height.SetText(strconv.Itoa(a.Preferences().IntWithFallback("canonical_height", 1632)))
	v.output = widget.NewMultiLineEntry()
	v.output.SetMinRowsVisible(3)
	v.output.SetPlaceHolder("Il messaggio verificato apparirà qui")
	v.output.Disable()
	v.file = wrap("Nessuna immagine selezionata · JPEG o PNG")
	v.preview = canvas.NewImageFromResource(theme.FileImageIcon())
	v.preview.FillMode = canvas.ImageFillContain
	v.preview.SetMinSize(fyne.NewSize(280, 100))
	v.preview.Hide()
	v.choose = widget.NewButtonWithIcon("Scegli immagine", theme.FolderOpenIcon(), v.openImage)
	v.read = widget.NewButtonWithIcon("Leggi il messaggio", theme.SearchIcon(), v.extract)
	v.read.Importance = widget.HighImportance
	v.create = widget.NewButtonWithIcon("Crea il PNG", theme.DocumentCreateIcon(), v.embed)
	v.create.Importance = widget.HighImportance
	v.mode = widget.NewSelect([]string{"Foto di una stampa", "PNG originale"}, func(s string) { v.updateControls() })
	v.mode.SetSelected("Foto di una stampa")
	dimensions := container.NewGridWithColumns(2, field("Larghezza (px)", v.width), field("Altezza (px)", v.height))
	readBody := container.NewVBox(field("Tipo di immagine", v.mode), dimensions, wrap("Per una foto usa le dimensioni del PNG prima della stampa. Inquadra l’immagine intera e i bordi del foglio."), v.read)
	createBody := container.NewVBox(field("Messaggio", v.message), v.byteCount, wrap("Prepara un PNG fino a 1632 pixel per le prove di stampa. Viene verificato prima del salvataggio."), v.create)
	v.tabs = container.NewAppTabs(container.NewTabItemWithIcon("Leggi", theme.SearchIcon(), readBody), container.NewTabItemWithIcon("Crea", theme.DocumentCreateIcon(), createBody))
	v.status = wrap("Pronto per una prova")
	v.status.TextStyle = fyne.TextStyle{Bold: true}
	v.elapsed = wrap("")
	v.progress = widget.NewProgressBarInfinite()
	v.progress.Hide()
	v.savePNG = widget.NewButtonWithIcon("Salva PNG", theme.DocumentSaveIcon(), func() { v.save(v.last.PNG, "pixseal.png") })
	v.saveReport = widget.NewButtonWithIcon("Salva resoconto", theme.DocumentSaveIcon(), v.exportReport)
	v.usePNG = widget.NewButton("Leggi questo PNG", v.useGenerated)
	v.copyMessage = widget.NewButtonWithIcon("Copia messaggio", theme.ContentCopyIcon(), func() {
		if !v.busy.Load() && v.last.Report.Authenticated && v.last.Report.Operation == "extract" {
			v.app.Clipboard().SetContent(v.output.Text)
		}
	})
	result := container.NewVBox(v.status, v.elapsed, v.progress, v.output, v.copyMessage, v.savePNG, v.usePNG, v.saveReport)
	title := canvas.NewText("PixSeal", pixTheme{}.Color(theme.ColorNamePrimary, theme.VariantLight))
	title.TextSize = 30
	title.TextStyle = fyne.TextStyle{Bold: true}
	body := container.NewVBox(title, wrap("Il messaggio viaggia nell’immagine."), widget.NewSeparator(), v.choose, v.file, container.NewCenter(v.preview), field("Chiave segreta", v.key), v.tabs, widget.NewSeparator(), result, wrap("Versione sperimentale · tutto viene elaborato sul dispositivo. Mantieni l’app aperta durante la lettura."))
	w.SetContent(container.NewVScroll(container.NewPadded(body)))
	v.updateControls()
	return v
}
func (v *View) updateControls() {
	if v.choose == nil {
		return
	}
	v.setControls(v.busy.Load())
}

func (v *View) setControls(busy bool) {
	for _, x := range []*widget.Button{v.choose, v.read, v.create} {
		if busy {
			x.Disable()
		} else {
			x.Enable()
		}
	}
	for _, x := range []*widget.Entry{v.key, v.message} {
		if busy {
			x.Disable()
		} else {
			x.Enable()
		}
	}
	if v.mode != nil {
		if busy {
			v.mode.Disable()
		} else {
			v.mode.Enable()
		}
	}
	for _, x := range []*widget.Entry{v.width, v.height} {
		if busy || (v.mode != nil && v.mode.Selected != "Foto di una stampa") {
			x.Disable()
		} else {
			x.Enable()
		}
	}
	for _, x := range []*widget.Button{v.savePNG, v.usePNG} {
		if x != nil {
			if busy || len(v.last.PNG) == 0 {
				x.Disable()
			} else {
				x.Enable()
			}
			if len(v.last.PNG) == 0 {
				x.Hide()
			} else {
				x.Show()
			}
		}
	}
	if v.saveReport != nil {
		if v.last.Report.Operation == "" {
			v.saveReport.Hide()
		} else {
			v.saveReport.Show()
		}
		if busy || v.last.Report.Operation == "" {
			v.saveReport.Disable()
		} else {
			v.saveReport.Enable()
		}
	}

	if v.copyMessage != nil {
		if v.last.Report.Authenticated && v.last.Report.Operation == "extract" {
			v.copyMessage.Show()
			v.output.Show()
		} else {
			v.copyMessage.Hide()
			v.output.Hide()
		}
		if busy || !v.last.Report.Authenticated || v.last.Report.Operation != "extract" {
			v.copyMessage.Disable()
		} else {
			v.copyMessage.Enable()
		}
	}
}

func (v *View) resetResult() { v.last = appservice.Result{}; v.output.SetText(""); v.updateControls() }
func (v *View) run(label string, clear bool, work func() (func(), error)) {
	if !v.busy.CompareAndSwap(false, true) {
		return
	}
	if clear {
		v.resetResult()
	}
	v.updateControls()
	v.status.SetText(label)
	v.progress.Show()
	v.progress.Start()
	v.elapsed.SetText("0 s")
	id := v.generation.Add(1)
	started := time.Now()
	finished := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-finished:
				return
			case <-ticker.C:
				fyne.Do(func() {
					if v.busy.Load() && v.generation.Load() == id {
						v.elapsed.SetText(fmt.Sprintf("%d s", int(time.Since(started).Seconds())))
					}
				})
			}
		}
	}()
	go func() {
		apply, err := work()
		close(finished)
		fyne.Do(func() {
			if apply != nil {
				apply()
			}
			if err != nil {
				v.status.SetText(err.Error())
			}
			v.progress.Stop()
			v.progress.Hide()
			v.elapsed.SetText(fmt.Sprintf("%.1f s", time.Since(started).Seconds()))
			v.setControls(false)
			v.busy.Store(false)
		})
	}()
}
func (v *View) setDocument(d appservice.Document) {
	v.doc = d
	v.file.SetText(fmt.Sprintf("%s · %d × %d px", d.Name, d.Image.Bounds().Dx(), d.Image.Bounds().Dy()))
	v.preview.Resource = nil
	v.preview.Image = appservice.Preview(d.Image)
	v.preview.Show()
	v.preview.Refresh()
}
func (v *View) openImage() {
	if v.busy.Load() {
		return
	}
	picker := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if r == nil {
			return
		}
		name := r.URI().Name()
		v.run("Apertura dell’immagine…", true, func() (func(), error) {
			defer r.Close()
			d, e := v.engine.Load(r, name)
			if e != nil {
				return nil, e
			}
			return func() { v.setDocument(d); v.status.SetText("Immagine pronta") }, nil
		})
	}, v.window)
	picker.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg", "image/png"}))
	picker.Show()
}
func (v *View) embed() {
	d, m, k := v.doc, v.message.Text, v.key.Text
	v.run("Creazione e verifica del PNG…", true, func() (func(), error) {
		r, err := v.engine.Embed(d, m, k)
		if err != nil {
			return nil, err
		}
		return func() {
			v.last = r
			v.width.SetText(strconv.Itoa(r.Report.CanonicalWidth))
			v.height.SetText(strconv.Itoa(r.Report.CanonicalHeight))
			v.app.Preferences().SetInt("canonical_width", r.Report.CanonicalWidth)
			v.app.Preferences().SetInt("canonical_height", r.Report.CanonicalHeight)
			v.status.SetText(fmt.Sprintf("PNG verificato · %d × %d px", r.Report.CanonicalWidth, r.Report.CanonicalHeight))
		}, nil
	})
}
func (v *View) extract() {
	d, k := v.doc, v.key.Text
	mode := "photo"
	if v.mode.Selected == "PNG originale" {
		mode = "original"
	}
	w, _ := strconv.Atoi(strings.TrimSpace(v.width.Text))
	h, _ := strconv.Atoi(strings.TrimSpace(v.height.Text))
	v.run("Ricerca del messaggio… può richiedere alcuni minuti", true, func() (func(), error) {
		r, err := v.engine.Extract(d, k, mode, w, h)
		if err == nil && !r.Report.Authenticated {
			err = appservice.ErrNotAuthenticated
		}
		if err != nil {
			r.Report.Authenticated = false
			r.Message = ""
		}
		return func() {
			v.last = appservice.Result{Report: r.Report}
			if err == nil && r.Report.Authenticated {
				v.output.SetText(r.Message)
				v.last = r
				v.status.SetText("Messaggio autenticato")
			}
		}, err
	})
}
func (v *View) useGenerated() {
	if v.busy.Load() || len(v.last.PNG) == 0 {
		return
	}
	hash := sha256.Sum256(v.last.PNG)
	v.setDocument(appservice.Document{Image: v.last.Image, Name: "pixseal.png", Format: "png", SHA256: hex.EncodeToString(hash[:])})
	v.mode.SetSelected("PNG originale")
	v.tabs.SelectIndex(0)
	v.status.SetText("PNG creato selezionato per la lettura")
}
func (v *View) exportReport() {
	if v.last.Report.Operation == "" {
		return
	}
	data, err := json.MarshalIndent(v.last.Report, "", "  ")
	if err != nil {
		dialog.ShowError(err, v.window)
		return
	}
	v.save(append(data, '\n'), "pixseal-prova.json")
}
func (v *View) save(data []byte, name string) {
	if v.busy.Load() || len(data) == 0 {
		return
	}
	picker := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, v.window)
			return
		}
		if w == nil {
			return
		}
		v.run("Salvataggio…", false, func() (func(), error) {
			n, writeErr := w.Write(data)
			closeErr := w.Close()
			if writeErr == nil && n != len(data) {
				writeErr = io.ErrShortWrite
			}
			err := errors.Join(writeErr, closeErr)
			if err != nil {
				return nil, err
			}
			return func() { v.status.SetText("File salvato") }, nil
		})
	}, v.window)
	picker.SetFileName(name)
	if strings.HasSuffix(name, ".png") {
		picker.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/png"}))
	} else {
		picker.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/json"}))
	}
	picker.Show()
}
