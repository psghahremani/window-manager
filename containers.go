package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"
import "fmt"

type Geometry struct {
	X int64
	Y int64

	Width  int64
	Height int64
}

type Container struct {
	parent   *Container
	children []*Container

	isSplitHorizontally bool

	containerWindow Window
	clientWindow    *Window
}

var VisibleWindows map[Window]*Container

func init() {
	VisibleWindows = map[Window]*Container{}
}

func CreateNewChildContainerWithNoParentToChildConnection(parentContainer *Container) *Container {
	childContainerWindow := NewSimpleWindow(parentContainer.containerWindow)
	childContainer := &Container{
		parent:              parentContainer,
		isSplitHorizontally: parentContainer.isSplitHorizontally,
		containerWindow:     (Window)(childContainerWindow),
	}

	childContainer.containerWindow.Map()
	VisibleWindows[childContainer.containerWindow] = childContainer

	//_ = childContainer.containerWindow.SubscribeToXEvents(
	//	[]C.long{
	//		C.SubstructureRedirectMask,
	//		C.SubstructureNotifyMask,
	//		//C.KeyPressMask,
	//	},
	//)

	fmt.Println("Grabbed key for window:", childContainer.containerWindow)

	//C.XGrabKey(
	//	xConnection,
	//	C.AnyKey,
	//	C.Mod4Mask,
	//	(C.Window)(childContainer.containerWindow),
	//	1,
	//	C.GrabModeSync,
	//	C.GrabModeSync,
	//)

	return childContainer
}

func (c *Container) CreateNewChildOrReuseExistingChildWithoutRefreshingGeometry(childIndex int64) *Container {
	if childIndex < (int64)(len(c.children)) {
		targetedContainer := c.children[childIndex]
		if targetedContainer.clientWindow == nil {
			// This container is an empty leaf container.
			return targetedContainer
		}

		// A new empty leaf container must be created and then shifted to the desired index.
		c.children = append(
			c.children,
			CreateNewChildContainerWithNoParentToChildConnection(c),
		)
		for i := (int64)(len(c.children) - 1); i > childIndex; i-- {
			c.children[i], c.children[i-1] = c.children[i-1], c.children[i]
		}
		return c.children[childIndex]
	}

	// Child index is out of range so empty "padding" containers are inserted and last one is returned.
	for i := (int64)(len(c.children)); i <= childIndex; i++ {
		c.children = append(
			c.children,
			CreateNewChildContainerWithNoParentToChildConnection(c),
		)
	}
	return c.children[len(c.children)-1]
}

func (c *Container) SetGeometry(geometry Geometry) {
	c.containerWindow.SetGeometry(geometry)

	if c.clientWindow != nil {
		c.clientWindow.SetGeometry(
			Geometry{
				X:      0,
				Y:      0,
				Width:  geometry.Width,
				Height: geometry.Height,
			},
		)
		return
	}

	if len(c.children) > 0 {
		if c.isSplitHorizontally {
			childrenY := (int64)(0)

			childrenHeight := geometry.Height / (int64)(len(c.children))
			residueHeight := geometry.Height % (int64)(len(c.children))

			for childIndex, child := range c.children {
				childHeight := childrenHeight
				if (int64)(childIndex) < residueHeight {
					childHeight += 1
				}

				childGeometry := Geometry{
					X:     geometry.X,
					Width: geometry.Width,

					Y:      childrenY,
					Height: childHeight,
				}
				child.SetGeometry(childGeometry)

				childrenY += childHeight
			}
		} else {
			childrenX := (int64)(0)

			childrenWidth := geometry.Width / (int64)(len(c.children))
			residueWidth := geometry.Width % (int64)(len(c.children))

			for childIndex, child := range c.children {
				childWidth := childrenWidth
				if (int64)(childIndex) < residueWidth {
					childWidth += 1
				}

				childGeometry := Geometry{
					Y:      geometry.Y,
					Height: geometry.Height,

					X:     childrenX,
					Width: childWidth,
				}
				child.SetGeometry(childGeometry)

				childrenX += childWidth
			}
		}
	}
}

func (c *Container) RefreshGeometry() {
	c.SetGeometry(c.containerWindow.GetGeometry())
}

func (c *Container) CreatePath(path []int64) *Container {
	currentContainer := c
	for _, index := range path {
		currentContainer = currentContainer.CreateNewChildOrReuseExistingChildWithoutRefreshingGeometry(index)
	}
	c.RefreshGeometry()
	return currentContainer
}

func (c *Container) SetClientWindow(clientWindow Window) {
	c.clientWindow = &clientWindow
	c.clientWindow.RemoveBorder()
	c.clientWindow.SetParentWindow(c.containerWindow)

	//_ = c.clientWindow.SubscribeToXEvents(
	//	[]C.long{
	//		C.SubstructureRedirectMask,
	//		C.SubstructureNotifyMask,
	//		//C.KeyPressMask,
	//	},
	//)

	geometry := c.containerWindow.GetGeometry()
	c.clientWindow.SetGeometry(
		Geometry{
			X:      0,
			Y:      0,
			Width:  geometry.Width,
			Height: geometry.Height,
		},
	)

	c.clientWindow.Map()
	VisibleWindows[*c.clientWindow] = c
}

func (c *Container) ToggleSplitDirection() {
	c.isSplitHorizontally = !c.isSplitHorizontally
	c.RefreshGeometry()
}

func (c *Container) HandleKeyEvent(modifierMask int64, keyCode int64) {
	killClientWindow := func() {
		// TODO: Use graceful termination on client windows.
		c.clientWindow.Kill()
		delete(VisibleWindows, (Window)(*c.clientWindow))
		c.clientWindow = nil
	}

	if keyCode == 24 {
		if len(c.children) > 0 {
			fmt.Println("Could not destroy container. It has children.")
			return
		}

		if c.clientWindow != nil {
			killClientWindow()
		}

		c.containerWindow.Kill()
		delete(VisibleWindows, c.containerWindow)

		for childIndexInParent, childInParent := range c.parent.children {
			if childInParent == c {
				c.parent.children = append(
					c.parent.children[:childIndexInParent],
					c.parent.children[childIndexInParent+1:]...,
				)
				break
			}
		}

		c.parent.RefreshGeometry()
	} else if keyCode == 54 {
		if c.clientWindow != nil {
			killClientWindow()
		}
	}
}
