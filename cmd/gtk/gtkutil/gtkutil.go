//go111:build gtk

package gtkutil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// updateMessageDialog makes MessageDialog labels selectable and markup-enabled.
// https://stackoverflow.com/questions/3249053/copying-the-text-from-a-gtk-messagedialog
func updateMessageDialog(dlg *gtk.MessageDialog) {
	// area, err := dlg.GetMessageArea()
	// if err == nil {
	// 	children := area.GetChildren()
	// 	children.Foreach(func(item any) {
	// 		label, err := gtk.WidgetToLabel(item.(*gtk.Widget))
	// 		if err == nil {
	// 			label.SetSelectable(true)
	// 			label.SetUseMarkup(true)
	// 		}
	// 	})
	// }
}

func ShowQuestionDialog(parent *gtk.Window, msg string) bool {
	return IdleAddSyncT(func() bool {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageQuestion, gtk.ButtonsYesNo)
		dlg.SetMarkup(msg)
		updateMessageDialog(dlg)

		responseChan := make(chan gtk.ResponseType, 1)
		dlg.Connect("response", func(responseID int) {
			responseChan <- gtk.ResponseType(responseID)
			dlg.Destroy()
		})
		dlg.SetVisible(true)
		response := <-responseChan

		return response == gtk.ResponseYes
	})
}

func ShowInfoDialog(parent *gtk.Window, msg string) {
	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageInfo, gtk.ButtonsOK)
		dlg.SetMarkup(msg)
		updateMessageDialog(dlg)
		ShowModalDialog(&dlg.Window)
	})
}

func ShowWarningDialog(parent *gtk.Window, msg string) {
	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageWarning, gtk.ButtonsOK)
		dlg.SetMarkup(msg)
		updateMessageDialog(dlg)
		ShowModalDialog(&dlg.Window)
	})
}

func ShowErrorDialog(parent *gtk.Window, msg string) {
	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageError, gtk.ButtonsOK)
		dlg.SetMarkup(msg)
		updateMessageDialog(dlg)
		ShowModalDialog(&dlg.Window)
	})
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

// ImageOption represents an option for ImageFromBytes.
type ImageOption func(*ImageOptions)

type ImageOptions struct {
	width  int
	height int
}

// WithImageSize sets the desired width and height for the image.
func WithImageSize(width, height int) ImageOption {
	return func(opts *ImageOptions) {
		opts.width = width
		opts.height = height
	}
}

// ImageFromBytes creates a gtk.Image from a byte slice and applies a marginEnd.
// It returns an empty gtk.Image if an error occurs.
func ImageFromBytes(data []byte, opts ...ImageOption) *gtk.Image {
	texture, err := gdk.NewTextureFromBytes(glib.NewBytes(data))
	if err != nil {
		Logf("Error creating texture from bytes: %v\n", err)

		return gtk.NewImage()
	}

	img := gtk.NewImageFromPaintable(&texture.Paintable)

	if len(opts) > 0 {
		options := &ImageOptions{}
		for _, opt := range opts {
			opt(options)
		}

		img = ResizeImage(img, options.width, options.height)
	}

	return img
}

func ResizeImage(img *gtk.Image, width, height int) *gtk.Image {
	img2 := gtk.NewImageFromPaintable(img.Paintable())
	img2.SetSizeRequest(int(width), int(height))

	return img2
}

func GetTextViewContent(tv *gtk.TextView) string {
	buf := tv.Buffer()
	startIter, endIter := buf.Bounds()
	content := buf.Text(startIter, endIter, true)

	return content
}

func SetTextViewContent(tv *gtk.TextView, content string) {
	buf := tv.Buffer()
	buf.SetText(content)
}

// OpenURLInBrowser opens a URL in the OS default browser.
func OpenURLInBrowser(address string) error {
	var cmd string
	args := make([]string, 0, 2)

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
	obj := builder.GetObject(overlayID)
	overlay := obj.Cast().(*gtk.Overlay)

	// Create a new Entry
	entry := gtk.NewEntry()
	entry.SetCanFocus(true)
	entry.SetHExpand(true)
	entry.SetEditable(false)

	SetCSSClass(&entry.Widget, "copyable_entry")

	// Create a new Button
	button := gtk.NewButtonFromIconName("edit-copy-symbolic")

	button.SetTooltipText("Copy to Clipboard")
	button.SetHAlign(gtk.AlignEnd)
	button.SetVAlign(gtk.AlignCenter)
	button.SetHExpand(false)
	button.SetVExpand(false)

	SetCSSClass(&button.Widget, "inline_button")

	// Set the click event for the Button
	button.Connect("clicked", func() {
		buffer := GetEntryText(entry)
		clipboard := button.Clipboard()
		clipboard.SetText(buffer)
	})

	overlay.SetChild(entry)
	overlay.AddOverlay(button)

	overlay.SetVisible(true)

	return entry
}

func SetCSSClass(widget *gtk.Widget, name string) {
	styleContext := widget.StyleContext()

	styleContext.AddClass(name)
}

func ShowNonModalDialog(dlg *gtk.Window) {
	IdleAddSync(func() {
		dlg.SetModal(false)
		dlg.SetVisible(true)
	})
}

func ShowModalDialog(dlg *gtk.Window) {
	IdleAddSync(func() {
		dlg.SetModal(true)
		dlg.SetVisible(true)
	})
}

func ComboBoxActiveValue(combo *gtk.ComboBox) int {
	iter, _ := combo.ActiveIter()
	model := combo.Model()
	val := model.Cast().(*gtk.TreeModel).Value(iter, 0)
	valueInterface := val.GoValue()

	return valueInterface.(int)
}

func GetEntryText(entry *gtk.Entry) string {
	return entry.Text()
}

func GetDropDown(drop *gtk.DropDown) string {
	return drop.SelectedItem().Type().String()
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

func IdleAddAsync(fun func()) {
	go func() {
		glib.IdleAdd(func() bool {
			fun()

			return false
		})
	}()
}

func IdleAddSync(fun func()) {
	IdleAddSyncT(func() bool {
		fun()

		return false
	})
}

func IdleAddSyncT[T any](fun func() T) T {
	res, _ := IdleAddSyncTT(func() (T, bool) {
		return fun(), false
	})

	return res
}

func IdleAddSyncTT[T1, T2 any](fun func() (T1, T2)) (T1, T2) {
	done := make(chan bool, 1)
	var va1l T1
	var val2 T2

	go func() {
		glib.IdleAdd(func() bool {
			va1l, val2 = fun()

			done <- true

			return false
		})
	}()

	glibContext := glib.MainContextDefault()
	for {
		select {
		case <-done:
			return va1l, val2
		default:
			time.Sleep(10 * time.Millisecond)
			glibContext.Iteration(false)
		}
	}
}

func Logf(msg string, args ...any) {
	log.Printf("(Go Routine ID %d) %s", GoroutineID(), fmt.Sprintf(msg, args...))
}

func GoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	var id int64
	_, _ = fmt.Sscanf(string(buf[:n]), "goroutine %d ", &id)

	return id
}

func AddImageToButton(button *gtk.Button, image *gtk.Image) {

	button.SetChild(image)
}

func AppendRowToListStore(listStore *gtk.ListStore, cols []int, values []any) {
	iter := listStore.Append()
	glibValues := make([]glib.Value, len(values))
	for i, v := range values {
		glibValues[i] = *glib.NewValue(v)
	}
	listStore.Set(iter, cols, glibValues)
}
