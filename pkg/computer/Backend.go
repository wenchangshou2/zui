package computer

type PointerButton int
type Key int

const (
	PointerButtonLeft PointerButton = iota
	PointerButtonRight
	PointerButtonMiddle
	PointerButtonLimit
)

const (
	KeyVolumeMute Key = iota
	KeyVolumeDown
	KeyVolumeUp
	KeyMediaPlayPause
	KeyMediaPrevTrack
	KeyMediaNextTrack
	KeyLimit
)

type BackendInfo struct {
	Name string
	Init func() (Backend, error)
}

var Backends []BackendInfo = []BackendInfo{
	{"Windows", InitWindowsBackend},
}

type UnsupportedPlatformError struct {
	err error
}

func (e UnsupportedPlatformError) Error() string {
	return e.err.Error()
}

type Backend interface {
	Close() error
	KeyboardText(text string) error
	KeyboardKey(key Key) error
	PointerButton(button PointerButton, press bool) error
	PointerMove(deltaX, deltaY int) error
	PointerScroll(deltaHorizontal, deltaVertical int) error
	PointerScrollFinish() error
	KeyTap(tapKey string, args ...interface{}) string
}



