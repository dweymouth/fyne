//go:build darwin

package glfw

/*
#import <stdbool.h>

void setDisableDisplaySleep(bool);
double doubleClickInterval();
void installFyneReopenDelegate(void);
*/
import "C"
import (
	"time"

	"fyne.io/fyne/v2"
)

func setDisableScreenBlank(disable bool) {
	C.setDisableDisplaySleep(C.bool(disable))
}

func (d *gLDriver) DoubleTapDelay() time.Duration {
	millis := int64(float64(C.doubleClickInterval()) * 1000)
	return time.Duration(millis) * time.Millisecond
}

func (d *gLDriver) initPlatform() {
	C.installFyneReopenDelegate()
}

//export fyneAppReopened
func fyneAppReopened() {
	// This is called from within glfw.PollEvents(), so the main goroutine is
	// blocked and cannot drain funcQueue. Dispatch to a new goroutine so the
	// callback returns immediately, unblocking PollEvents, before we call
	// runOnMain to show the window.
	go func() {
		d := fyne.CurrentApp().Driver().(*gLDriver)
		runOnMain(func() {
			for _, w := range d.windows {
				win := w.(*window)
				if win.master {
					win.Show()
					return
				}
			}
		})
	}()
}
