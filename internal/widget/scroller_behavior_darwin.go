//go:build !ci && darwin

package widget

/*
#cgo LDFLAGS: -framework AppKit

int getScrollerPagingBehavior();
int getScrollerStyle();
void watchScrollerStyle();
*/
import "C"

import (
	"sync"

	"fyne.io/fyne/v2"
)

func isScrollerPageOnTap() bool {
	return C.getScrollerPagingBehavior() == 0
}

// scrollBarAlwaysVisible returns true when the macOS "Show scroll bars" preference
// is set to "Always" (NSScrollerStyleLegacy = 0).
// When set to "When scrolling" or "Automatically" (NSScrollerStyleOverlay = 1),
// scroll bars should only appear on hover or while scrolling.
func scrollBarAlwaysVisible() bool {
	return C.getScrollerStyle() == 0
}

var (
	scrollerStyleWatchOnce   sync.Once
	scrollerStyleMu          sync.Mutex
	scrollerStyleSubscribers = map[uint64]func(){}
	scrollerStyleNextID      uint64
)

func subscribeScrollerStyle(fn func()) uint64 {
	scrollerStyleWatchOnce.Do(func() { C.watchScrollerStyle() })
	scrollerStyleMu.Lock()
	defer scrollerStyleMu.Unlock()
	id := scrollerStyleNextID
	scrollerStyleNextID++
	scrollerStyleSubscribers[id] = fn
	return id
}

func unsubscribeScrollerStyle(id uint64) {
	scrollerStyleMu.Lock()
	defer scrollerStyleMu.Unlock()
	delete(scrollerStyleSubscribers, id)
}

//export scrollerStyleChanged
func scrollerStyleChanged() {
	scrollerStyleMu.Lock()
	fns := make([]func(), 0, len(scrollerStyleSubscribers))
	for _, fn := range scrollerStyleSubscribers {
		fns = append(fns, fn)
	}
	scrollerStyleMu.Unlock()
	fyne.Do(func() {
		for _, fn := range fns {
			fn()
		}
	})
}
