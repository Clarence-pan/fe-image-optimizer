package main

//#include "load-resource_windows.h"
import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unsafe"
)

func loadAppResourceById(resId int) (startAddr *byte, size int) {
	handle := C.loadAppResourceById(
		C.int64(resId),
		(**C.byte)(unsafe.Pointer(&startAddr)),
		(*C.int64)(unsafe.Pointer(&size)),
	)

	if handle == nil {
		panic(fmt.Errorf("failed to load resource by id(%d), lastError: %#v", resId, getLastError()))
	}

	return
}

func getLastError() error {
	// lastErr := C.getLastError()
	lastErr := C.getLastErrorAsString()
	return fmt.Errorf("#%s", C.GoString(lastErr))
}

func extractAppResource(resId int, filename string, mode os.FileMode) {
	addr, size := loadAppResourceById(resId)
	log.Printf("extract res#%d to %s (size=%v)", resId, filename, size)

	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = *((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) + uintptr(i))))
	}

	mkdirIfNotExists(filepath.Dir(filename))

	ensure(ioutil.WriteFile(filename, buf, mode))
}
