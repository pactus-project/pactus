//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

// SplashWindow is a lightweight splash screen that stays visible while the
// node boots up. Callers can update the status text via SetStatus.
type SplashWindow struct {
	window      *gtk.ApplicationWindow
	statusLabel *gtk.Label
	version     *gtk.Label
	spinner     *gtk.Spinner
}

func NewSplashWindow(app *gtk.Application) *SplashWindow {
	window, err := gtk.ApplicationWindowNew(app)
	gtkutil.FatalErrorCheck(err)

	window.SetDecorated(false)
	window.SetResizable(false)
	window.SetTypeHint(gdk.WINDOW_TYPE_HINT_SPLASHSCREEN)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetDefaultSize(420, 220)

	content, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)
	gtkutil.FatalErrorCheck(err)
	content.SetBorderWidth(24)

	logo := gtkutil.ImageFromPixbuf(assets.ImagePactusLogoPixbuf)
	if logo != nil {
		logo.SetMarginBottom(8)
		content.Add(logo)
	}

	spinner, err := gtk.SpinnerNew()
	gtkutil.FatalErrorCheck(err)
	spinner.Start()
	spinner.SetMarginBottom(6)
	content.Add(spinner)

	statusLabel, err := gtk.LabelNew("Starting node...")
	gtkutil.FatalErrorCheck(err)
	statusLabel.SetHAlign(gtk.ALIGN_START)
	content.Add(statusLabel)

	versionLabel, err := gtk.LabelNew("")
	gtkutil.FatalErrorCheck(err)
	versionLabel.SetHAlign(gtk.ALIGN_START)
	versionLabel.SetMarginTop(6)
	styleContext, err := versionLabel.GetStyleContext()
	gtkutil.FatalErrorCheck(err)
	styleContext.AddClass("dim-label")
	content.Add(versionLabel)

	window.Add(content)

	return &SplashWindow{
		window:      window,
		statusLabel: statusLabel,
		spinner:     spinner,
		version:     versionLabel,
	}
}

func (s *SplashWindow) ShowAll() {
	s.window.ShowAll()
	s.spinner.Start()
}

func (s *SplashWindow) Destroy() {
	s.window.Destroy()
}

func (s *SplashWindow) SetStatus(text string) {
	s.statusLabel.SetText(text)
}

func (s *SplashWindow) SetVersion(text string) {
	s.version.SetText(text)
}

func (s *SplashWindow) Window() *gtk.ApplicationWindow {
	return s.window
}
