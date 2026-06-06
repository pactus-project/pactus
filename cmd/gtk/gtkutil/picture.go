//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func NewScaledPictureFromTexture(texture *gdk.Texture, width, height int) *gtk.Picture {
	picture := gtk.NewPicture()
	picture.SetPaintable(texture)
	picture.SetSizeRequest(width, height)
	picture.SetHAlign(gtk.AlignCenter)
	picture.SetVAlign(gtk.AlignCenter)

	picture.SetContentFit(gtk.ContentFitContain)

	return picture
}
