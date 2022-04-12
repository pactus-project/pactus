package main

import (
	"flag"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/util"
)

var (
	workingDir *string
)

func init() {
	workingDir = flag.String("working-dir", cmd.ZarbHomeDir(), "working directory")
}
func main() {
	flag.Parse()

	gtk.Init(nil)

	if util.IsDirNotExistsOrEmpty(*workingDir) {
		if !startupWizard() {
			
			return
		}
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Zarb GUI")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	//
}
