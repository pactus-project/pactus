//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

// SplashWindow is a lightweight splash screen that stays visible while the
// node is starting.
type SplashWindow struct {
	window      *gtk.ApplicationWindow
	statusLabel *gtk.Label
	version     *gtk.Label
	spinner     *gtk.Spinner
}

func NewSplashWindow(app *gtk.Application) *SplashWindow {
	window := gtk.NewApplicationWindow(app)

	// Configure window appearance
	window.SetDecorated(false)
	window.SetResizable(false)
	window.SetDefaultSize(420, 220)

	// Create main container with proper padding and spacing
	mainBox := gtk.NewBox(gtk.OrientationVertical, 20)
	mainBox.AddCSSClass("splash")

	// Logo section
	pic := gtkutil.NewScaledPictureFromTexture(assets.ImagePactusLogoTexture, 128, 128)
	pic.AddCSSClass("splash-logo")
	mainBox.Append(pic)

	// Spinner and status section
	spinner := gtk.NewSpinner()
	pic.AddCSSClass("splash-spinner")
	spinner.Start()
	spinner.SetHAlign(gtk.AlignCenter)
	spinner.SetSizeRequest(12, 12)
	mainBox.Append(spinner)

	statusLabel := gtk.NewLabel("Starting node...")
	statusLabel.SetHAlign(gtk.AlignStart)
	statusLabel.AddCSSClass("splash-status")
	mainBox.Append(statusLabel)

	// Version label at bottom
	versionLabel := gtk.NewLabel("")
	versionLabel.SetHAlign(gtk.AlignStart)
	versionLabel.AddCSSClass("splash-version")
	mainBox.Append(versionLabel)

	window.SetChild(mainBox)

	return &SplashWindow{
		window:      window,
		statusLabel: statusLabel,
		spinner:     spinner,
		version:     versionLabel,
	}
}

func (s *SplashWindow) ShowAll() {
	s.window.Present()
	s.spinner.Start()

	// The splash has no parent to center over and GTK4 cannot move a window,
	// so center it natively once it has been mapped.
	glib.TimeoutAdd(80, func() bool {
		gtkutil.CenterActiveWindow()

		return false
	})
}

func (s *SplashWindow) Destroy() {
	s.window.Destroy()
}

func (s *SplashWindow) UpdateStatus(text string) {
	s.statusLabel.SetText(text)
}

func (s *SplashWindow) SetVersion(text string) {
	s.version.SetText(text)
}

func (s *SplashWindow) Window() *gtk.Window {
	return &s.window.Window
}
