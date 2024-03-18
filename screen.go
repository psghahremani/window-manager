package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

type Screen C.Screen

func (s *Screen) GetDimensions() (int64, int64) {
	cScreen := (*C.Screen)(s)
	return (int64)(cScreen.width), (int64)(cScreen.height)
}

func (s *Screen) GetRootWindow() Window {
	cScreen := (*C.Screen)(s)
	cRootWindow := C.XRootWindowOfScreen(cScreen)
	return (Window)(cRootWindow)
}
