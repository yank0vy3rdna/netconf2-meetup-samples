//nolint:nlreturn
package sysrepo

/*
#include "helper.h"
*/
import "C"

import (
	"unsafe"
)

type data C.sr_data_t

func (d *data) C() *C.sr_data_t {
	return (*C.sr_data_t)(d)
}

func (d *data) DataTreeToString(format LydFormat) string {
	var str *C.char

	C.lyd_print_mem(&str, d.tree, format.C(), LYD_PRINT_WITHSIBLINGS.C())
	defer C.free(unsafe.Pointer(str))

	return C.GoString(str)
}

func (d *data) Free() {
	C.sr_release_data(d.C())
}
