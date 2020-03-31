package windows
import "C"



func ActivePID(pid int32,args ...int)error{
	var hwnd int
	if len(args)>0{
		hwnd=args[0]
	}
	internalActive(pid,hwnd)
	return nil
}

func internalActive(pid int32, hwnd interface{}) {
	C.active_PID(C.uintptr(pid),C.uinptr(hwnd))
}