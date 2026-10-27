package main

import (
	"errors"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/pj/pixseal/internal/appservice"
)

type stubEngine struct {
	gate   chan struct{}
	result appservice.Result
	err    error
	calls  chan string
}

func (s *stubEngine) Load(io.Reader, string) (appservice.Document, error) {
	return appservice.Document{}, errors.New("unused")
}
func (s *stubEngine) Embed(appservice.Document, string, string) (appservice.Result, error) {
	return s.result, s.err
}
func (s *stubEngine) Extract(_ appservice.Document, _ string, mode string, _, _ int) (appservice.Result, error) {
	if s.calls != nil {
		s.calls <- mode
	}
	if s.gate != nil {
		<-s.gate
	}
	return s.result, s.err
}
func waitView(t *testing.T, v *View) {
	t.Helper()
	deadline := time.Now().Add(4 * time.Second)
	for v.busy.Load() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if v.busy.Load() {
		t.Fatal("GUI worker did not finish")
	}
}
func TestReadHidesStaleAndUnauthenticatedMessages(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("PixSeal")
	e := &stubEngine{gate: make(chan struct{}), calls: make(chan string, 1), result: appservice.Result{Message: "UNTRUSTED", Report: appservice.Report{Operation: "extract", Authenticated: true}}, err: errors.New("non autenticato")}
	v := NewView(a, w, e)
	v.doc = appservice.Document{Image: image.NewNRGBA(image.Rect(0, 0, 592, 512))}
	v.key.SetText("test-key-105")
	v.output.SetText("old successful message")
	test.Tap(v.read)
	if <-e.calls != "photo" {
		t.Fatal("wrong mode")
	}
	if v.output.Text != "" || !v.read.Disabled() || !v.create.Disabled() || !v.choose.Disabled() {
		t.Fatal("UI exposes stale data or permits concurrent jobs")
	}
	close(e.gate)
	waitView(t, v)
	if v.output.Text != "" || !v.copyMessage.Disabled() || v.last.Report.Authenticated {
		t.Fatal("UI exposed unauthenticated message")
	}
	e.gate = nil
	e.err = nil
	e.result = appservice.Result{Message: "Verified", Report: appservice.Report{Operation: "extract", Authenticated: true}}
	v.mode.SetSelected("PNG originale")
	test.Tap(v.read)
	if <-e.calls != "original" {
		t.Fatal("wrong original mode")
	}
	waitView(t, v)
	if v.output.Text != "Verified" || v.copyMessage.Disabled() {
		t.Fatalf("verified message unavailable: output=%q status=%q copyDisabled=%t", v.output.Text, v.status.Text, v.copyMessage.Disabled())
	}
	test.Tap(v.copyMessage)
	if a.Clipboard().Content() != "Verified" {
		t.Fatal("copy failed")
	}
}
func TestUIRequiresAuthenticationEvenWithoutBackendError(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	w := a.NewWindow("PixSeal")
	v := NewView(a, w, &stubEngine{result: appservice.Result{Message: "DO NOT SHOW", Report: appservice.Report{Operation: "extract"}}})
	test.Tap(v.read)
	waitView(t, v)
	if v.output.Text != "" || v.status.Text == "Messaggio autenticato" {
		t.Fatal("non-authenticated payload displayed")
	}
}
func TestMobileLayouts(t *testing.T) {
	// ci builds use Fyne's real software painter, enabling a native-widget capture.
	a := app.NewWithID("io.github.piggei.pixseal.tests")
	defer a.Quit()
	a.Settings().SetTheme(pixTheme{})
	w := a.NewWindow("PixSeal")
	v := NewView(a, w, &stubEngine{})
	w.Resize(fyne.NewSize(390, 820))
	w.Show()
	dir := os.Getenv("PIXSEAL_GUI_SNAPSHOT_DIR")
	for i, name := range []string{"read", "create"} {
		v.tabs.SelectIndex(i)
		if w.Content().MinSize().Width > 390 {
			t.Fatal("mobile content overflows")
		}
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			f, err := os.Create(filepath.Join(dir, "pixseal-android-"+name+".png"))
			if err != nil {
				t.Fatal(err)
			}
			err = png.Encode(f, w.Canvas().Capture())
			closeErr := f.Close()
			if err != nil || closeErr != nil {
				t.Fatal(err, closeErr)
			}
		}
	}
}
