
package main

/*
#cgo pkg-config: x11
#include <X11/Xlib.h>
*/
import "C"

import "unsafe"

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
        event_ptr := unsafe.Pointer(&event)

        switch *(*C.int)(event_ptr) {

        case C.KeyPress:

            xkey := *(*C.XKeyEvent)(event_ptr)

            if xkey.subwindow == C.None { break }
            C.XRaiseWindow(dpy, xkey.subwindow)

        case C.ButtonPress:

            xbutton := *(*C.XButtonEvent)(event_ptr)

            if xbutton.subwindow == C.None { break }
            C.XGetWindowAttributes(dpy, xbutton.subwindow, &attr)
            start = xbutton

        case C.MotionNotify:

            xbutton := *(*C.XButtonEvent)(event_ptr)

            if start.subwindow == C.None { break }
            
            xdiff := xbutton.x_root - start.x_root
            ydiff := xbutton.y_root - start.y_root

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