//go:build gtk

package gtkutil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// UpdateMessageDialog makes MessageDialog labels selectable and markup-enabled.
// https://stackoverflow.com/questions/3249053/copying-the-text-from-a-gtk-messagedialog
func UpdateMessageDialog(dlg *gtk.MessageDialog) {
	area, err := dlg.GetMessageArea()
	if err == nil {
		children := area.GetChildren()
		children.Foreach(func(item any) {
			label, err := gtk.WidgetToLabel(item.(*gtk.Widget))
			if err == nil {
				label.SetSelectable(true)
				label.SetUseMarkup(true)
			}
		})
	}
}

func ShowQuestionDialog(parent gtk.IWindow, msg string) bool {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "%s", msg)
	UpdateMessageDialog(dlg)
	res := RunDialog(&dlg.Dialog)

	return res == gtk.RESPONSE_YES
}

func ShowInfoDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", msg)
	UpdateMessageDialog(dlg)
	RunDialog(&dlg.Dialog)
}

func ShowWarningDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", msg)
	UpdateMessageDialog(dlg)
	RunDialog(&dlg.Dialog)
}

func ShowErrorDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "%s", msg)
	UpdateMessageDialog(dlg)
	RunDialog(&dlg.Dialog)
}

// ShowError displays an error dialog and logs the error message.
func ShowError(err error) {
	ShowErrorDialog(nil, err.Error())
	log.Print(err.Error())
}

// FatalErrorCheck checks for an error, shows an error dialog and terminates the program.
// Use with caution.
func FatalErrorCheck(err error) {
	if err != nil {
		ShowErrorDialog(nil, err.Error())
		log.Fatal(err.Error())
	}
}

// PixbufOption represents an option for PixbufFromBytes.
type PixbufOption func(*pixbufOptions)

type pixbufOptions struct {
	width  int
	height int
}

// WithSize sets the desired width and height for the pixbuf.
func WithSize(width, height int) PixbufOption {
	return func(opts *pixbufOptions) {
		opts.width = width
		opts.height = height
	}
}

// PixbufFromBytes decodes image bytes (PNG/SVG/etc) into a gdk.Pixbuf.
// Use WithSize() option to resize the image.
func PixbufFromBytes(data []byte, opts ...PixbufOption) (*gdk.Pixbuf, error) {
	options := &pixbufOptions{}
	for _, opt := range opts {
		opt(options)
	}

	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = loader.Close()
		if err != nil {
			log.Println("error closing pixbuf loader:", err)
		}
	}()

	pixbuf, err := loader.WriteAndReturnPixbuf(data)
	if err != nil {
		return nil, err
	}

	// Resize if size options provided
	if options.width > 0 && options.height > 0 {
		resized, err := pixbuf.ScaleSimple(options.width, options.height, gdk.INTERP_NEAREST)
		if err != nil {
			return nil, err
		}

		return resized, nil
	}

	return pixbuf, nil
}

// ImageFromPixbuf creates a gtk.Image from a pixbuf and applies a marginEnd.
// If pixbuf is nil, it returns an empty gtk.Image.
func ImageFromPixbuf(pixbuf *gdk.Pixbuf) *gtk.Image {
	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil
	}

	image.ShowAll()

	return image
}

func GetTextViewContent(tv *gtk.TextView) string {
	buf, _ := tv.GetBuffer()
	startIter, endIter := buf.GetBounds()
	content, err := buf.GetText(startIter, endIter, true)
	if err != nil {
		return ""
	}

	return content
}

func SetTextViewContent(tv *gtk.TextView, content string) {
	buf, err := tv.GetBuffer()
	if err != nil {
		return
	}
	buf.SetText(content)
}

// OpenURLInBrowser opens a URL in the OS default browser.
func OpenURLInBrowser(address string) error {
	var cmd string
	args := make([]string, 0)

	addr, err := url.Parse(address)
	if err != nil {
		return err
	}

	switch addr.Scheme {
	case "http", "https":
	default:
		return errors.New("address scheme is invalid")
	}

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, address)

	return exec.CommandContext(context.Background(), cmd, args...).Start()
}

func BuildExtendedEntry(builder *gtk.Builder, overlayID string) *gtk.Entry {
	obj, err := builder.GetObject(overlayID)
	FatalErrorCheck(err)
	overlay := obj.(*gtk.Overlay)

	// Create a new Entry
	entry, err := gtk.EntryNew()
	FatalErrorCheck(err)
	entry.SetCanFocus(true)
	entry.SetHExpand(true)
	entry.SetEditable(false)

	SetCSSClass(&entry.Widget, "copyable_entry")

	// Create a new Button
	button, err := gtk.ButtonNewFromIconName("edit-copy-symbolic", gtk.ICON_SIZE_BUTTON)
	FatalErrorCheck(err)

	button.SetTooltipText("Copy to Clipboard") // TODO: Not working!
	button.SetHAlign(gtk.ALIGN_END)
	button.SetVAlign(gtk.ALIGN_CENTER)
	button.SetHExpand(false)
	button.SetVExpand(false)
	button.SetBorderWidth(0)

	SetCSSClass(&button.Widget, "inline_button")

	// Set the click event for the Button
	button.Connect("clicked", func() {
		buffer := GetEntryText(entry)
		clipboard, _ := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
		clipboard.SetText(buffer)
	})

	overlay.Add(entry)
	overlay.AddOverlay(button)

	overlay.ShowAll() // Ensure all child widgets are shown

	return entry
}

func SetCSSClass(widget *gtk.Widget, name string) {
	styleContext, err := widget.GetStyleContext()
	FatalErrorCheck(err)

	styleContext.AddClass(name)
}

func RunDialog(dlg *gtk.Dialog) gtk.ResponseType {
	response := dlg.Run()

	// Destroy should be done after the dialog is closed
	// Read more here: https://docs.gtk.org/gtk3/method.Dialog.run.html
	dlg.Destroy()

	return response
}

func ComboBoxActiveValue(combo *gtk.ComboBox) int {
	iter, err := combo.GetActiveIter()
	FatalErrorCheck(err)

	model, err := combo.GetModel()
	FatalErrorCheck(err)

	val, err := model.ToTreeModel().GetValue(iter, 0)
	FatalErrorCheck(err)

	valueInterface, err := val.GoValue()
	FatalErrorCheck(err)

	return valueInterface.(int)
}

func GetEntryText(entry *gtk.Entry) string {
	txt, err := entry.GetText()
	FatalErrorCheck(err)

	return txt
}

// Color represents different text colors for UI elements.
type Color int

const (
	ColorRed Color = iota
	ColorGreen
	ColorBlue
	ColorYellow
	ColorOrange
	ColorPurple
	ColorGray
)

// getColorHex returns the hex color code for the given Color enum.
func getColorHex(color Color) string {
	switch color {
	case ColorRed:
		return "#FF0000"
	case ColorGreen:
		return "#00FF00"
	case ColorBlue:
		return "#0000FF"
	case ColorYellow:
		return "#FFFF00"
	case ColorOrange:
		return "#FFA500"
	case ColorPurple:
		return "#800080"
	case ColorGray:
		return "#808080"
	default:
		return "#000000" // Default to black
	}
}

func SetColoredText(label *gtk.Label, str string, color Color) {
	colorHex := getColorHex(color)
	formattedText := fmt.Sprintf("<span color='%s'>%s</span>", colorHex, str)
	label.SetMarkup(formattedText)
}

func IdleAddAsync(f func()) {
	go func() {
		glib.IdleAdd(func() bool {
			f()

			return false
		})
	}()
}

func IdleAddSync(f func()) {
	IdleAddSyncT(func() bool {
		f()

		return false
	})
}

func IdleAddSyncT[T any](f func() T) T {
	res, _ := IdleAddSyncTT(func() (T, bool) {
		return f(), false
	})

	return res
}

func IdleAddSyncTT[T1, T2 any](f func() (T1, T2)) (T1, T2) {
	type pair struct {
		A T1
		B T2
	}

	done := make(chan bool)
	var a T1
	var b T2

	go func() {
		fmt.Print("\n.")
		glib.IdleAdd(func() bool {
			fmt.Print("-")
			a, b = f()
			fmt.Print("+")

			done <- true

			return false
		})
	}()
	<-done

	return a, b
}
