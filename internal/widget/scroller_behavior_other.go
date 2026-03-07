//go:build !darwin

package widget

func isScrollerPageOnTap() bool {
	return false
}

func scrollBarAlwaysVisible() bool {
	return true
}

func subscribeScrollerStyle(_ func()) int { return 0 }
func unsubscribeScrollerStyle(_ int)      {}
