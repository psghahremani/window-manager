package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

import (
	"errors"
)

var xConnection *C.Display

func ConnectToXServer() error {
	display := C.XOpenDisplay(nil)
	if display == nil {
		return errors.New("could not connect to X server")
	}
	xConnection = display

	C.XSynchronize(xConnection, 1)
	return nil
}

func LockServer() {
	C.XGrabServer(xConnection)
}

func UnlockServer() {
	C.XUngrabServer(xConnection)
}

func GetMainScreen() (*Screen, error) {
	cScreen := C.XScreenOfDisplay(xConnection, 0)
	return (*Screen)(cScreen), nil
}

func Flush() {
	C.XFlush(xConnection)
}
