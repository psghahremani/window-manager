package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <stdlib.h>
import "C"

import (
	"reflect"
	"unsafe"
)

type Window C.Window

type WindowAttributes struct {
	IsVisible                    bool
	SubstructureRedirectOverride bool
	DoNotPropagateMask           int64
}

func NewSimpleWindow(parentWindow Window) Window {
	window := C.XCreateSimpleWindow(
		xConnection,
		(C.Window)(parentWindow),
		0,
		0,
		1,
		1,
		0,
		0x000000,
		0x000000,
	)
	return (Window)(window)
}

func (w Window) SetParentWindow(parentWindow Window) {
	C.XReparentWindow(xConnection, (C.Window)(w), (C.Window)(parentWindow), 0, 0)
}

func (w Window) GetWindowAttributes() WindowAttributes {
	var windowAttributes C.XWindowAttributes
	C.XGetWindowAttributes(xConnection, (C.Window)(w), &windowAttributes)

	return WindowAttributes{
		IsVisible:                    windowAttributes.map_state == 2,
		SubstructureRedirectOverride: windowAttributes.override_redirect != 0,
		DoNotPropagateMask:           (int64)(windowAttributes.do_not_propagate_mask),
	}
}

func (w Window) GetParentAndChildren() (Window, []Window) {
	var root, parent Window
	var childrenArrayPointer *Window
	var childrenCount C.uint
	C.XQueryTree(
		xConnection,
		(C.Window)(w),
		(*C.Window)(&root),
		(*C.Window)(&parent),
		(**C.Window)((unsafe.Pointer)(&childrenArrayPointer)),
		&childrenCount,
	)

	var children []Window
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&children))
	sliceHeader.Data = (uintptr)((unsafe.Pointer)(childrenArrayPointer))
	sliceHeader.Len = (int)(childrenCount)
	sliceHeader.Cap = (int)(childrenCount)

	// TODO: Is this still needed?
	// C.XFree((unsafe.Pointer)(childrenArrayPointer))

	return parent, children
}

func (w Window) SetGeometry(geometry Geometry) {
	var changes C.XWindowChanges
	changes.x = (C.int)(geometry.X)
	changes.y = (C.int)(geometry.Y)
	changes.width = (C.int)(geometry.Width)
	changes.height = (C.int)(geometry.Height)

	var mask int64 = C.CWX | C.CWY | C.CWWidth | C.CWHeight

	C.XConfigureWindow(
		xConnection,
		(C.Window)(w),
		(C.uint)(mask),
		&changes,
	)
}

func (w Window) GetGeometry() Geometry {
	var windowAttributes C.XWindowAttributes
	C.XGetWindowAttributes(xConnection, (C.Window)(w), &windowAttributes)

	return Geometry{
		X:      (int64)(windowAttributes.x),
		Y:      (int64)(windowAttributes.y),
		Width:  (int64)(windowAttributes.width),
		Height: (int64)(windowAttributes.height),
	}
}

func (w Window) SubscribeToXEvents(subscriptionMasks []C.long) error {
	finalMask := (C.long)(0)
	for _, mask := range subscriptionMasks {
		finalMask |= mask
	}

	//finalMask = C.NoEventMask |
	//	C.KeyPressMask |
	//	C.KeyReleaseMask |
	//	C.ButtonPressMask |
	//	C.ButtonReleaseMask |
	//	C.EnterWindowMask |
	//	C.LeaveWindowMask |
	//	C.PointerMotionMask |
	//	C.PointerMotionHintMask |
	//	C.Button1MotionMask |
	//	C.Button2MotionMask |
	//	C.Button3MotionMask |
	//	C.Button4MotionMask |
	//	C.Button5MotionMask |
	//	C.ButtonMotionMask |
	//	C.KeymapStateMask |
	//	C.ExposureMask |
	//	C.VisibilityChangeMask |
	//	C.StructureNotifyMask |
	//	C.ResizeRedirectMask |
	//	C.SubstructureNotifyMask |
	//	C.SubstructureRedirectMask |
	//	C.FocusChangeMask |
	//	C.PropertyChangeMask |
	//	C.ColormapChangeMask |
	//	C.OwnerGrabButtonMask

	C.XSelectInput(xConnection, (C.Window)(w), (C.long)(finalMask))

	return nil
}

func (w Window) Map() {
	C.XMapWindow(xConnection, (C.Window)(w))
}

func (w Window) Unmap() {
	C.XUnmapWindow(xConnection, (C.Window)(w))
}

func (w Window) Kill() {
	C.XDestroyWindow(xConnection, (C.Window)(w))
}

func (w Window) Clear() {
	C.XClearWindow(xConnection, (C.Window)(w))
}

func (w Window) RemoveBorder() {
	var changes C.XWindowChanges
	changes.border_width = 0

	var mask int64 = C.CWBorderWidth

	C.XConfigureWindow(
		xConnection,
		(C.Window)(w),
		(C.uint)(mask),
		&changes,
	)
}

func (w Window) GetName() string {
	var windowName *C.char
	C.XFetchName(xConnection, (C.Window)(w), &windowName)
	defer C.free((unsafe.Pointer)(windowName))
	return C.GoString(windowName)
}
