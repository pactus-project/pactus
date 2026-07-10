//go:build gtk

package view

import (
	"fmt"
	"math"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Pactus accent teal (#12909e), matching the CSS accent color.
const (
	accentRed   = 0x12 / 255.0
	accentGreen = 0x90 / 255.0
	accentBlue  = 0x9e / 255.0
)

// CircularProgress is a ring gauge that fills clockwise with the sync progress
// and shows the percentage in its center. The track color follows the current
// theme so it reads well in both light and dark mode.
type CircularProgress struct {
	*gtk.Overlay

	area     *gtk.DrawingArea
	label    *gtk.Label
	fraction float64
}

// NewCircularProgress creates a ring gauge with the given square size in pixels.
func NewCircularProgress(size int) *CircularProgress {
	area := gtk.NewDrawingArea()
	area.SetContentWidth(size)
	area.SetContentHeight(size)

	label := gtk.NewLabel("0%")
	label.SetHAlign(gtk.AlignCenter)
	label.SetVAlign(gtk.AlignCenter)
	label.AddCSSClass("circular-progress-label")

	overlay := gtk.NewOverlay()
	overlay.SetChild(area)
	overlay.AddOverlay(label)

	progress := &CircularProgress{
		Overlay: overlay,
		area:    area,
		label:   label,
	}

	area.SetDrawFunc(func(_ *gtk.DrawingArea, cr *cairo.Context, width, height int) {
		progress.draw(cr, width, height)
	})

	return progress
}

// SetFraction sets the progress in the range [0, 1] and refreshes the ring.
func (cp *CircularProgress) SetFraction(fraction float64) {
	fraction = math.Max(0, math.Min(1, fraction))
	cp.fraction = fraction
	cp.label.SetText(fmt.Sprintf("%.0f%%", fraction*100))
	cp.area.QueueDraw()
}

func (cp *CircularProgress) draw(cr *cairo.Context, width, height int) {
	w := float64(width)
	h := float64(height)
	lineWidth := math.Max(6, math.Min(w, h)*0.09)
	radius := math.Min(w, h)/2 - lineWidth
	centerX := w / 2
	centerY := h / 2

	cr.SetLineWidth(lineWidth)
	cr.SetLineCap(cairo.LineCapRound)

	// Track: widget foreground color at low opacity, so it adapts to the theme.
	fg := cp.area.Color()
	cr.SetSourceRGBA(float64(fg.Red()), float64(fg.Green()), float64(fg.Blue()), 0.15)
	cr.Arc(centerX, centerY, radius, 0, 2*math.Pi)
	cr.Stroke()

	// Progress arc, clockwise starting from the top (-90 degrees).
	if cp.fraction > 0 {
		start := -math.Pi / 2
		end := start + 2*math.Pi*cp.fraction
		cr.SetSourceRGBA(accentRed, accentGreen, accentBlue, 1)
		cr.Arc(centerX, centerY, radius, start, end)
		cr.Stroke()
	}
}
