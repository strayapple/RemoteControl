package mouse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/vblazhnov/RemoteControl/interfaces"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	setCursorPosInfo = user32.NewProc("SetCursorPos")
	getCursorPosInfo = user32.NewProc("GetCursorPos")
	mouseEventInfo   = user32.NewProc("mouse_event")
)

const (
	clickTypeLeft int = iota
	clickTypeRight
)

// Control allow to remote use mouse
type Control struct {
}

type position struct {
	x, y int32
}

// GetHandlers implemintation of interfaces.EndPointInfo
func (c Control) GetHandlers() []interfaces.EndPointInfo {
	handlers := []interfaces.EndPointInfo{
		interfaces.BaseEndPointInfo{
			"/mouse/move",
			moveServerRequest,
			interfaces.EndPointAPI},
		interfaces.BaseEndPointInfo{
			"/mouse/click",
			clickServerRequest,
			interfaces.EndPointAPI},
		interfaces.BaseEndPointInfo{
			"/mouse",
			contentServeRequest,
			interfaces.EndPointContent},
	}

	return handlers
}

func contentServeRequest(
	w http.ResponseWriter,
	r *http.Request) {
	f, e := ioutil.ReadFile(".\\plugins\\mouse\\static\\index.html")
	if e != nil {
		fmt.Fprint(w, e)
	}
	w.Write(f)
}

func moveServerRequest(
	w http.ResponseWriter,
	r *http.Request) {
	xStr, yStr := r.FormValue("x"), r.FormValue("y")

	x, errX := strconv.ParseInt(xStr, 0, 32)
	y, errY := strconv.ParseInt(yStr, 0, 32)
	if errX != nil || errY != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	currentPos := position{}
	getCursorPosInfo.Call(uintptr(unsafe.Pointer(&currentPos)))
	setCursorPosInfo.Call(
		uintptr(currentPos.x+int32(x)),
		uintptr(currentPos.y+int32(y)),
	)
	fmt.Fprint(w, currentPos.x+int32(x), currentPos.y+int32(y))
}

func clickServerRequest(
	w http.ResponseWriter,
	r *http.Request) {
	typeStr := r.FormValue("type")
	clickType, err := strconv.Atoi(typeStr)
	if err != nil || clickType > clickTypeRight {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	click(clickType)
}

const (
	mouseEventFLeftDown  uint32 = 0x02
	mouseEventFLeftUp    uint32 = 0x04
	mouseEventFRightDown uint32 = 0x08
	mouseEventFRightUp   uint32 = 0x10
)

func click(t int) {
	switch t {
	case clickTypeLeft:
		mouseEventInfo.Call(
			uintptr(mouseEventFLeftDown|mouseEventFLeftUp),
			uintptr(uint(0)),
			uintptr(uint(0)),
			uintptr(uint(0)),
			uintptr(uint(0)),
		)

	case clickTypeRight:
		mouseEventInfo.Call(
			uintptr(mouseEventFRightDown|mouseEventFRightUp),
			uintptr(uint(0)),
			uintptr(uint(0)),
			uintptr(uint(0)),
			uintptr(uint(0)),
		)
	}
}