package main
//#cgo CFLAGS: -std=gnu99 -I/usr/include
//#include "sshpass.h"
import "C"

func sshpass(host, password string) {
	host_c := C.CString(host)
	password_c := C.CString(password)

	C.sshpass(host_c, password_c)
}
