//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// The GUI draws its icons from bundled SVG files instead of emoji, because
// emoji depend on the fonts of the operating system and may show up as empty
// boxes (for example a bold warning sign, U+26A0, on Windows).

const (
	iconSizeSmall  = 16
	iconSizeDialog = 48
)

// NewIcon returns an image that shows the texture at the given pixel size.
func NewIcon(texture *gdk.Texture, size int) *gtk.Image {
	img := gtk.NewImageFromPaintable(texture)
	img.SetPixelSize(size)

	return img
}

// NewIconButton returns a button that shows an icon in front of its label.
func NewIconButton(texture *gdk.Texture, label string) *gtk.Button {
	box := gtk.NewBox(gtk.OrientationHorizontal, 6)
	box.SetHAlign(gtk.AlignCenter)
	box.Append(NewIcon(texture, iconSizeSmall))
	box.Append(gtk.NewLabel(label))

	button := gtk.NewButton()
	button.SetChild(box)

	return button
}

// StatusLabel is a label with an optional icon in front of it, used for status
// messages that change over time.
type StatusLabel struct {
	*gtk.Box

	Icon  *gtk.Image
	Label *gtk.Label
}

func NewStatusLabel() *StatusLabel {
	icon := gtk.NewImage()
	icon.SetPixelSize(iconSizeSmall)
	icon.SetVAlign(gtk.AlignStart)
	icon.SetVisible(false)

	label := gtk.NewLabel("")
	label.SetXAlign(0)

	box := gtk.NewBox(gtk.OrientationHorizontal, 6)
	box.Append(icon)
	box.Append(label)

	return &StatusLabel{Box: box, Icon: icon, Label: label}
}

// SetText shows the text with the icon in front. A nil icon hides the icon.
func (s *StatusLabel) SetText(icon *gdk.Texture, text string) {
	s.setIcon(icon)
	s.Label.SetText(text)
}

// SetColoredText shows the text in the given color with the icon in front.
func (s *StatusLabel) SetColoredText(icon *gdk.Texture, text string, color Color) {
	s.setIcon(icon)
	SetColoredText(s.Label, text, color)
}

// Clear removes the icon and the text.
func (s *StatusLabel) Clear() {
	s.setIcon(nil)
	s.Label.SetText("")
}

func (s *StatusLabel) setIcon(icon *gdk.Texture) {
	if icon == nil {
		s.Icon.Clear()
		s.Icon.SetVisible(false)

		return
	}
	s.Icon.SetFromPaintable(icon)
	s.Icon.SetVisible(true)
}
