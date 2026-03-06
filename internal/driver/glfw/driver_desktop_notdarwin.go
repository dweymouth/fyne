//go:build !darwin && !wasm && !test_web_driver

package glfw

func (d *gLDriver) initPlatform() {}
