//package main
//
//import "C"
//import (
//	"fmt"
//	"time"
//	"unsafe"
//)
//
//// #cgo LDFLAGS: -lX11
//// #include <X11/Xlib.h>
//// #include <X11/keysymdef.h>
//// typedef void (*closure)();
//// int handleXError(Display *display, XErrorEvent *event);
//import "C"
//
////export handleXError
//func handleXError(display *C.Display, event *C.XErrorEvent) C.int {
//	fmt.Println("Error occurred during XGrabKey")
//	return 0
//}
//
//func main() {
//	C.XSetErrorHandler(C.closure(C.handleXError))
//
//	err := ConnectToXServer()
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	mainScreen, err := GetMainScreen()
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	mainScreenWidth, mainScreenHeight := mainScreen.GetDimensions()
//	fmt.Printf("Main screen is %d by %d pixels wide.\n", mainScreenWidth, mainScreenHeight)
//	rootWindow := mainScreen.GetRootWindow()
//
//	// TODO:
//	// Root container and root window can not have a 1:1 mapping.
//	// We need workspaces. So the mapping would be (1 x RootWindow) <-> (N x RootContainer).
//	// Switching between workspaces is simply a map/unmap for the root container.
//	rootContainer := &Container{
//		isSplitHorizontally: false,
//		containerWindow:     rootWindow,
//	}
//	VisibleWindows[rootWindow] = rootContainer
//
//	err = rootContainer.containerWindow.SubscribeToXEvents(
//		[]C.long{
//			C.SubstructureRedirectMask,
//			C.SubstructureNotifyMask,
//			//C.KeyPressMask,
//		},
//	)
//	if err != nil {
//		fmt.Println("Could not subscribe to events on root window. Is another window manager currently running?")
//		return
//	}
//
//	C.XGrabKey(
//		xConnection,
//		C.AnyKey,
//		0,
//		(C.Window)(rootWindow),
//		C.False,
//		C.GrabModeAsync,
//		C.GrabModeAsync,
//	)
//
//	LockServer()
//
//	_, existingClientWindows := rootWindow.GetParentAndChildren()
//	fmt.Printf("All Windows: %d\n", len(existingClientWindows))
//	visibleWindowsCount := 0
//	for _, clientWindow := range existingClientWindows {
//		clientWindowAttributes := clientWindow.GetWindowAttributes()
//		if clientWindowAttributes.IsVisible {
//			visibleWindowsCount += 1
//			container := rootContainer.CreatePath([]int64{0})
//			container.SetClientWindow(clientWindow)
//		} else {
//			fmt.Printf("An existing window %d,%s is invisible. Ignoring...\n", (int64)(clientWindow), clientWindow.GetName())
//		}
//	}
//	fmt.Printf("Visible Windows: %d\n", visibleWindowsCount)
//
//	printContainers(0, rootContainer)
//	fmt.Println("Done.")
//
//	UnlockServer()
//
//	for {
//		var event C.XEvent
//		C.XNextEvent(xConnection, &event)
//
//		eventType := (*C.int)(unsafe.Pointer(&event[0]))
//		fmt.Printf("Got Event: %d, %s\n", int64(*eventType), EventsMap[int64(*eventType)])
//
//		switch *eventType {
//		// Requests:
//		case C.MapNotify:
//			clientWindow := (Window)((*C.XMapEvent)(unsafe.Pointer(&event[0])).window)
//			parent, _ := clientWindow.GetParentAndChildren()
//			fmt.Printf(
//				"MapNotify received for a Window: (%d, %s) which is the child of %d, and override is set to %t and DonotPropagateMask is %d.\n",
//				(int64)(clientWindow),
//				clientWindow.GetName(),
//				(int64)(parent),
//				clientWindow.GetWindowAttributes().SubstructureRedirectOverride,
//				clientWindow.GetWindowAttributes().DoNotPropagateMask,
//			)
//		case C.MapRequest:
//			clientWindow := (Window)((*C.XMapRequestEvent)(unsafe.Pointer(&event[0])).window)
//			rootContainer.CreatePath([]int64{0}).SetClientWindow(clientWindow)
//			fmt.Printf("Mapped a Window: %d, %s, do_not: %d\n", (int64)(clientWindow), clientWindow.GetName(), clientWindow.GetWindowAttributes().DoNotPropagateMask)
//			printContainers(0, rootContainer)
//		case C.ConfigureRequest:
//			// TODO: Check if we're managing this window.
//
//			request := (*C.XConfigureRequestEvent)(unsafe.Pointer(&event[0]))
//
//			var notification C.XConfigureEvent
//			notification._type = C.ConfigureNotify
//			notification.display = xConnection
//			notification.event = request.window
//			notification.window = request.window
//
//			windowGeometry := (Window)(request.window).GetGeometry()
//			notification.x = (C.int)(windowGeometry.X)
//			notification.y = (C.int)(windowGeometry.Y)
//			notification.width = (C.int)(windowGeometry.Width)
//			notification.height = (C.int)(windowGeometry.Height)
//
//			notification.above = 0
//			notification.override_redirect = 0
//
//			var notificationEvent C.XEvent
//			*(*C.XConfigureEvent)(unsafe.Pointer(&notificationEvent[0])) = notification
//			C.XSendEvent(xConnection, request.window, 0, C.StructureNotifyMask, &notificationEvent)
//
//		// Notifications:
//		case C.ButtonPress:
//			fmt.Printf("Mouse button was pressed: %d\n", (*C.XButtonEvent)(unsafe.Pointer(&event[0])).button)
//		case C.KeyPress:
//			eventPayload := (*C.XKeyEvent)(unsafe.Pointer(&event[0]))
//			currentWindow := (Window)(eventPayload.window)
//			fmt.Printf("KeyStroke:\n")
//
//			if eventPayload.state&C.Mod4Mask == 0 {
//				fmt.Printf("Not related to WM (no Mod4).\n")
//				continue
//			}
//
//			fmt.Printf("Related to WM (has Mod4).\n")
//			fmt.Printf("Sleeping...\n")
//			time.Sleep(5 * time.Second)
//			fmt.Printf("Woke up...\n")
//
//			for {
//				fmt.Printf("Window:%d\n", currentWindow)
//				parent, _ := currentWindow.GetParentAndChildren()
//				currentWindow = parent
//				if currentWindow == 0 {
//					break
//				}
//			}
//			continue
//		//case C.KeyPress:
//		//	eventPayload := (*C.XKeyEvent)(unsafe.Pointer(&event[0]))
//		//	if eventPayload.state&C.Mod4Mask == 0 {
//		//		fmt.Println("PASS: Key event did not have Mod4.")
//		//
//		//		C.XAllowEvents(xConnection, C.ReplayKeyboard, eventPayload.time)
//		//		continue
//		//	}
//		//
//		//	fmt.Printf("Received key event root:%d,window:%d,subwindow:%d\n", eventPayload.root, eventPayload.window, eventPayload.subwindow)
//		//	fmt.Printf("Received key event: Mod:%d, KeyCode:%d\n", eventPayload.state, eventPayload.keycode)
//		//
//		//	if eventPayload.keycode == 36 {
//		//		err := exec.Command("xeyes").Start()
//		//		if err != nil {
//		//			fmt.Println("Could not launch terminal:", err.Error())
//		//			continue
//		//		}
//		//		fmt.Println("Launched terminal.")
//		//		continue
//		//	}
//		//
//		//	if eventPayload.window == (C.Window)(rootContainer.containerWindow) {
//		//		if eventPayload.subwindow == 0 {
//		//			fmt.Println("PASS: Key event was on an empty space.")
//		//
//		//			C.XAllowEvents(xConnection, C.ReplayKeyboard, eventPayload.time)
//		//			continue
//		//		}
//		//
//		//		container, isVisible := VisibleWindows[(Window)(eventPayload.subwindow)]
//		//		if !isVisible {
//		//			fmt.Println("PASS: Key event was ignored because it was directed to an invisible child of root.")
//		//
//		//			C.XAllowEvents(xConnection, C.ReplayKeyboard, eventPayload.time)
//		//			continue
//		//		}
//		//
//		//		fmt.Println("GRAB: Key event was sent to a container (root's direct child).")
//		//		container.HandleKeyEvent((int64)(eventPayload.state), (int64)(eventPayload.keycode))
//		//
//		//		C.XAllowEvents(xConnection, C.AsyncKeyboard, eventPayload.time)
//		//		continue
//		//	}
//		//
//		//	if eventPayload.subwindow != 0 {
//		//		container, isVisible := VisibleWindows[(Window)(eventPayload.subwindow)]
//		//		if !isVisible {
//		//			fmt.Println("PASS: IDK WTF JUST HAPPENED.")
//		//
//		//			C.XAllowEvents(xConnection, C.ReplayKeyboard, eventPayload.time)
//		//			continue
//		//		}
//		//
//		//		fmt.Println("GRAB: Key event was sent to a container (But it had a child container so the child got it).")
//		//		container.HandleKeyEvent((int64)(eventPayload.state), (int64)(eventPayload.keycode))
//		//
//		//		C.XAllowEvents(xConnection, C.AsyncKeyboard, eventPayload.time)
//		//		continue
//		//	}
//		//
//		//	container, isVisible := VisibleWindows[(Window)(eventPayload.window)]
//		//	if !isVisible {
//		//		fmt.Println("PASS: IDK WTF JUST HAPPENED.")
//		//
//		//		C.XAllowEvents(xConnection, C.ReplayKeyboard, eventPayload.time)
//		//		continue
//		//	}
//		//
//		//	fmt.Println("GRAB: Key event was sent to a container (root's indirect child).")
//		//	container.HandleKeyEvent((int64)(eventPayload.state), (int64)(eventPayload.keycode))
//		//
//		//	C.XAllowEvents(xConnection, C.AsyncKeyboard, eventPayload.time)
//		//	continue
//		case C.UnmapNotify:
//			window := (Window)((*C.XUnmapEvent)(unsafe.Pointer(&event[0])).window)
//			container, isVisible := VisibleWindows[window]
//			if !isVisible {
//				fmt.Println("Ignored because of invisibility.")
//				continue
//			}
//			if container.clientWindow != nil && *container.clientWindow == window {
//				fmt.Printf("Client window %d,%s has been unmapped. Container will be updated.\n", (int64)(*container.clientWindow), container.clientWindow.GetName())
//				container.clientWindow = nil
//				delete(VisibleWindows, window)
//				printContainers(0, rootContainer)
//			}
//		case C.DestroyNotify:
//			window := (Window)((*C.XDestroyWindowEvent)(unsafe.Pointer(&event[0])).window)
//			container, isVisible := VisibleWindows[window]
//			if !isVisible {
//				fmt.Println("Ignored because of invisibility.")
//				continue
//			}
//			if container.clientWindow != nil && *container.clientWindow == window {
//				fmt.Printf("Client window %d,%s has been destroyed. Container will be updated.\n", (int64)(*container.clientWindow), container.clientWindow.GetName())
//				container.clientWindow = nil
//				delete(VisibleWindows, window)
//				printContainers(0, rootContainer)
//			}
//		}
//	}
//}
//
//func printContainers(indent int64, container *Container) {
//	indentString := ""
//	for i := (int64)(0); i < indent; i++ {
//		indentString += "      "
//	}
//	indentString += "| "
//
//	fmt.Print(indentString)
//	fmt.Print("==========\n")
//
//	fmt.Print(indentString)
//	fmt.Printf("Position: (%d,%d)\n", container.containerWindow.GetGeometry().X, container.containerWindow.GetGeometry().Y)
//
//	fmt.Print(indentString)
//	fmt.Printf("Size: (%d,%d)\n", container.containerWindow.GetGeometry().Width, container.containerWindow.GetGeometry().Height)
//
//	fmt.Print(indentString)
//	fmt.Printf("ContainerWindow: %d\n", (int64)(container.containerWindow))
//
//	if container.clientWindow != nil {
//		fmt.Print(indentString)
//		fmt.Printf("ClientWindow: %d, %s\n", (int64)(*container.clientWindow), container.clientWindow.GetName())
//	}
//
//	fmt.Print(indentString)
//	fmt.Printf("IsSplitHorizontally: %t\n", container.isSplitHorizontally)
//
//	fmt.Print(indentString)
//	fmt.Printf("ChildrenCount: %d\n", len(container.children))
//
//	for _, child := range container.children {
//		printContainers(indent+1, child)
//	}
//
//	fmt.Print(indentString)
//	fmt.Print("==========\n")
//}

package main

import "fmt"

func main() {
	fmt.Println('\b')
}
