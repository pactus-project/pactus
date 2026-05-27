//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

// SplashWindow is a lightweight splash screen that stays visible while the
// node boots up. Callers can update the status text via SetStatus and
// progress text via SetProgress.
type SplashWindow struct {
	window        *gtk.ApplicationWindow
	statusLabel   *gtk.Label
	progressLabel *gtk.Label
	version       *gtk.Label
	spinner       *gtk.Spinner
	progressBar   *gtk.ProgressBar
}

func NewSplashWindow(app *gtk.Application) *SplashWindow {
	window := gtk.NewApplicationWindow(app)

	// Configure window appearance
	window.SetDecorated(false)
	window.SetResizable(false)
	window.SetDefaultSize(500, 380)

	// Create main container with proper padding and spacing
	mainBox := gtk.NewBox(gtk.OrientationVertical, 20)
	mainBox.SetHomogeneous(false)
	mainBox.SetMarginStart(32)
	mainBox.SetMarginEnd(32)
	mainBox.SetMarginTop(32)
	mainBox.SetMarginBottom(32)

	// Logo section
	logoBox := gtk.NewBox(gtk.OrientationVertical, 0)
	logoBox.SetHomogeneous(false)
	logoBox.SetHAlign(gtk.AlignCenter)

	logo := gtkutil.ResizeImage(assets.ImagePactusLogo, 140, 140)
	logoBox.Append(logo)
	mainBox.Append(logoBox)

	// Spinner and status section
	statusBox := gtk.NewBox(gtk.OrientationVertical, 12)
	statusBox.SetHomogeneous(false)

	spinner := gtk.NewSpinner()
	spinner.Start()
	spinner.SetHAlign(gtk.AlignCenter)
	spinner.SetSizeRequest(32, 32)
	statusBox.Append(spinner)

	statusLabel := gtk.NewLabel("Starting node...")
	statusLabel.SetHAlign(gtk.AlignCenter)
	statusLabel.SetWrap(true)
	statusLabel.SetJustify(gtk.JustifyCenter)
	statusLabel.AddCSSClass("splash-status")
	statusBox.Append(statusLabel)

	mainBox.Append(statusBox)

	// Progress section
	progressBox := gtk.NewBox(gtk.OrientationVertical, 8)
	progressBox.SetHomogeneous(false)

	progressBar := gtk.NewProgressBar()
	progressBar.SetShowText(true)
	progressBar.SetFraction(0.0)
	progressBar.AddCSSClass("splash-progress")
	progressBox.Append(progressBar)

	progressLabel := gtk.NewLabel("")
	progressLabel.SetHAlign(gtk.AlignCenter)
	progressLabel.SetWrap(true)
	progressLabel.SetJustify(gtk.JustifyCenter)
	progressLabel.AddCSSClass("splash-progress-text")
	progressBox.Append(progressLabel)

	mainBox.Append(progressBox)

	// Version label at bottom
	versionLabel := gtk.NewLabel("")
	versionLabel.SetHAlign(gtk.AlignCenter)
	versionLabel.SetMarginTop(8)
	versionLabel.AddCSSClass("dim-label")
	versionLabel.AddCSSClass("splash-version")
	mainBox.Append(versionLabel)

	// Apply styling
	applySplashWindowStyles(window)

	window.SetChild(mainBox)

	return &SplashWindow{
		window:        window,
		statusLabel:   statusLabel,
		progressLabel: progressLabel,
		spinner:       spinner,
		version:       versionLabel,
		progressBar:   progressBar,
	}
}

// applySplashWindowStyles applies CSS styling to the splash window.
func applySplashWindowStyles(window *gtk.ApplicationWindow) {
	provider := gtk.NewCSSProvider()
	css := `
		.splash-status {
			font-size: 14px;
			font-weight: 500;
			color: @theme_fg_color;
		}

		.splash-progress {
			min-height: 8px;
			border-radius: 4px;
		}

		.splash-progress-text {
			font-size: 11px;
			color: @theme_fg_color;
			opacity: 0.8;
		}

		.splash-version {
			font-size: 12px;
			opacity: 0.7;
		}
	`
	provider.LoadFromData(css)
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (s *SplashWindow) ShowAll() {
	s.window.Show()
	s.spinner.Start()
}

func (s *SplashWindow) Destroy() {
	s.window.Destroy()
}

func (s *SplashWindow) UpdateStatus(text string, progress int) {
	s.statusLabel.SetText(text)
	s.SetProgress(progress, "")
}

func (s *SplashWindow) SetProgress(fraction int, text string) {
	if fraction < 0 {
		return
	}

	s.progressBar.SetFraction(float64(fraction) / 100.0)
	s.progressLabel.SetText(text)
}

func (s *SplashWindow) SetVersion(text string) {
	s.version.SetText(text)
}

func (s *SplashWindow) Window() *gtk.Window {
	return &s.window.Window
}
