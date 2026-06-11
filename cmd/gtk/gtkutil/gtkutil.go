//go:build gtk

//nolint:staticcheck // Using depreciated widgets
package gtkutil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

func ShowQuestionDialog(parent *gtk.Window, msg string,
	onClone func(res gtk.ResponseType),
) {
	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageQuestion, gtk.ButtonsYesNo)

		showMessageDialog(dlg, "Question", msg, onClone)
	})
}

func ShowInfoDialog(parent *gtk.Window, msg string, onClone func(res gtk.ResponseType)) {
	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageInfo, gtk.ButtonsOK)
		showMessageDialog(dlg, "Info", msg, onClone)
	})
}

func ShowWarningDialog(parent *gtk.Window, msg string, onClone func(res gtk.ResponseType)) {
	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageWarning, gtk.ButtonsOK)
		showMessageDialog(dlg, "Warning", msg, onClone)
	})
}

func ShowErrorDialog(parent *gtk.Window, msg string, onClone func(res gtk.ResponseType)) {
	Logf("an error occurred: %s", msg)

	IdleAddSync(func() {
		dlg := gtk.NewMessageDialog(parent,
			gtk.DialogModal, gtk.MessageError, gtk.ButtonsOK)
		showMessageDialog(dlg, "Error", msg, onClone)
	})
}

func showMessageDialog(dlg *gtk.MessageDialog, title, msg string, onClose func(res gtk.ResponseType)) {
	dlg.SetMarkup(fmt.Sprintf("<b>%s</b>", title))
	dlg.SetObjectProperty("secondary-use-markup", true)
	dlg.SetObjectProperty("secondary-text", msg)

	dlg.ConnectResponse(func(responseID int) {
		dlg.Destroy()

		if onClose != nil {
			onClose(gtk.ResponseType(responseID))
		}
	})

	ShowModalWindow(&dlg.Window)
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
	overlay := builder.GetObject(overlayID).Cast().(*gtk.Overlay)

	// Create a new Entry
	entry := gtk.NewEntry()
	entry.SetCanFocus(true)
	entry.SetHExpand(true)
	entry.SetEditable(false)
	entry.AddCSSClass("copyable_entry")

	// Create a new Button
	button := gtk.NewButtonFromIconName("edit-copy-symbolic")

	button.SetTooltipText("Copy to Clipboard")
	button.SetHAlign(gtk.AlignEnd)
	button.SetVAlign(gtk.AlignCenter)
	button.SetHExpand(false)
	button.SetVExpand(false)
	button.AddCSSClass("inline_button")

	// Set the click event for the Button
	button.Connect("clicked", func() {
		buffer := EntryGetText(entry)
		clipboard := button.Clipboard()
		clipboard.SetText(buffer)
	})

	overlay.SetChild(entry)
	overlay.AddOverlay(button)

	overlay.SetVisible(true)

	return entry
}

func ShowNonModalWindow(win *gtk.Window) {
	IdleAddSync(func() {
		win.SetModal(false)
		win.Present()
	})
}

func ShowModalWindow(win *gtk.Window) {
	IdleAddAsync(func() {
		win.SetModal(true)
		win.Present()
	})
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

func UpdateSendButton(button *gtk.Button) {
	ExtendImageButton(button, "_Send", "Send this transaction", assets.IconSendTexture)
}

func UpdateOKButton(window *gtk.Window, button *gtk.Button) {
	ExtendImageButton(button, "_OK", "Perform this operation", assets.IconOkTexture)
	window.SetDefaultWidget(button)
}

func UpdateCancelButton(button *gtk.Button) {
	ExtendImageButton(button, "_Cancel", "Cancel this operation", assets.IconCancelTexture)
}

func UpdateCloseButton(button *gtk.Button) {
	ExtendImageButton(button, "_Close", "Close this window", assets.IconCloseTexture)
}

func ExtendImageButton(btn *gtk.Button, text, tooltip string, texture *gdk.Texture) {
	box := gtk.NewBox(gtk.OrientationHorizontal, 4)
	pic := NewScaledPictureFromTexture(texture, 16, 16)
	label := gtk.NewLabel(text)
	label.SetUseUnderline(true)

	box.Append(pic)
	box.Append(label)
	btn.SetChild(box)
	btn.SetTooltipText(tooltip)
}

func ConnectButtonSignal(button *gtk.Button, handler func()) {
	button.ConnectClicked(handler)
}

type ContextMenuItem struct {
	Label      string
	Action     func()
	IconPixbuf *gdkpixbuf.Pixbuf
}

func (c *ContextMenuItem) detailedAction() string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(strings.TrimSpace(c.Label)),
		"_", ""),
		" ", "-")
}

func CreateContextMenu(appWindow *gtk.ApplicationWindow, widget *gtk.Widget, menuItems []ContextMenuItem) {
	actionGroup := gio.NewSimpleActionGroup()

	for _, item := range menuItems {
		action := gio.NewSimpleAction(item.detailedAction(), nil)
		action.ConnectActivate(func(*glib.Variant) {
			item.Action()
		})
		action.SetEnabled(true)
		actionGroup.AddAction(action)
	}

	menuModel := gio.NewMenu()
	for _, item := range menuItems {
		menuModel.Append(item.Label, "win."+item.detailedAction())
	}
	appWindow.InsertActionGroup("win", actionGroup)

	popover := gtk.NewPopoverMenuFromModel(menuModel)
	popover.SetParent(widget)

	gesture := gtk.NewGestureClick()
	gesture.SetButton(3) // Right button

	gesture.ConnectPressed(func(_ int, x, y float64) {
		rect := gdk.NewRectangle(int(x), int(y), 1, 1)
		popover.SetPointingTo(&rect)
		popover.Popup()
	})

	widget.AddController(gesture)
}

func CaptureDoubleClick(widget *gtk.Widget, action func()) {
	gesture := gtk.NewGestureClick()
	gesture.SetButton(gdk.BUTTON_PRIMARY)
	gesture.ConnectPressed(func(nPress int, x, y float64) {
		if nPress == 2 { // Double-click
			Logf("Double-click detected at (%.2f, %.2f)", x, y)
			action()
		}
	})

	widget.AddController(gesture)
}

// ClearListModel removes all items from a gioutil ListModel.
func ClearListModel[T any](listModel *gioutil.ListModel[T]) {
	for i := int(listModel.NItems()); i > 0; i-- {
		listModel.Remove(i - 1)
	}
}

func IsWidgetShowing(widget *gtk.Widget) bool {
	return widget.Mapped()
}
