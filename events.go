package main

var EventsMap map[int64]string

func init() {
	EventsMap = map[int64]string{
		2:  "KeyPress",
		3:  "KeyRelease",
		4:  "ButtonPress",
		5:  "ButtonRelease",
		6:  "MotionNotify",
		7:  "EnterNotify",
		8:  "LeaveNotify",
		9:  "FocusIn",
		10: "FocusOut",
		11: "KeymapNotify",
		12: "Expose",
		13: "GraphicsExpose",
		14: "NoExpose",
		15: "VisibilityNotify",
		16: "CreateNotify",
		17: "DestroyNotify",
		18: "UnmapNotify",
		19: "MapNotify",
		20: "MapRequest",
		21: "ReparentNotify",
		22: "ConfigureNotify",
		23: "ConfigureRequest",
		24: "GravityNotify",
		25: "ResizeRequest",
		26: "CirculateNotify",
		27: "CirculateRequest",
		28: "PropertyNotify",
		29: "SelectionClear",
		30: "SelectionRequest",
		31: "SelectionNotify",
		32: "ColormapNotify",
		33: "ClientMessage",
		34: "MappingNotify",
		35: "GenericEvent",
		36: "LASTEvent",
	}
}
