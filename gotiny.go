
package main

/*
#cgo pkg-config: x11
#include <X11/Xlib.h>
int typeOf(XEvent event){ return event.type; }
XKeyEvent xkey(XEvent event) { return event.xkey; }
XButtonEvent xbutton(XEvent event) { return event.xbutton; }
*/
import "C"

func max(a, b int) int{
    if a < b { return b }
    return a
}

func main(){

    var start C.XButtonEvent
    var attr  C.XWindowAttributes
    var event C.XEvent;

    dpy := C.XOpenDisplay(nil)
    if dpy == nil { return }

    root_window := C.XDefaultRootWindow(dpy)
    
    C.XGrabKey(dpy, C.int(C.XKeysymToKeycode(dpy, root_window)), C.Mod1Mask, root_window,
               1, C.GrabModeAsync, C.GrabModeAsync)
    C.XGrabButton(dpy, 1, C.Mod1Mask, root_window, 1,
                  C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask, C.GrabModeAsync, C.GrabModeAsync, C.None, C.None);
    C.XGrabButton(dpy, 3, C.Mod1Mask, root_window, 1,
                  C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask, C.GrabModeAsync, C.GrabModeAsync, C.None, C.None);
    
    start.subwindow = C.None

    for{
        C.XNextEvent(dpy, &event)

        switch C.typeOf(event) {

        case C.KeyPress:
            if C.xkey(event).subwindow == C.None { break }
            C.XRaiseWindow(dpy, C.xkey(event).subwindow)

        case C.ButtonPress:
            if C.xbutton(event).subwindow == C.None { break }
            C.XGetWindowAttributes(dpy, C.xbutton(event).subwindow, &attr)
            start = C.xbutton(event)

        case C.MotionNotify:
            if start.subwindow == C.None { break }
            
            xdiff := C.xbutton(event).x_root - start.x_root
            ydiff := C.xbutton(event).y_root - start.y_root

            switch start.button {
            case 1:
                C.XMoveWindow(dpy, start.subwindow, attr.x + xdiff, attr.y + ydiff)
            case 3:
                C.XResizeWindow(dpy, start.subwindow,
                                C.uint(max(1, int(attr.width  + xdiff))),
                                C.uint(max(1, int(attr.height + ydiff))))
            }

        case C.ButtonRelease:
            start.subwindow = C.None
        }   
    }
}