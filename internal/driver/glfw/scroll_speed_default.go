//go:build !darwin && !wasm && !test_web_driver

package glfw

const (
	scrollAccelerateRate   = float64(200)
	scrollAccelerateCutoff = float64(10)
	scrollSpeed            = float32(75)
)
