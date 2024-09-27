package sysrepo

/*
#include "helper.h"
*/
import "C"

import (
	"errors"
	"fmt"
)

var (
	ErrItemNotFound = errors.New("item not found")
	ErrSysrepo      = errors.New("sysrepo error")
)

type Error int

const (
	ERR_OK                Error = C.SR_ERR_OK
	ERR_INVAL_ARG         Error = C.SR_ERR_INVAL_ARG
	ERR_LY                Error = C.SR_ERR_LY
	ERR_SYS               Error = C.SR_ERR_SYS
	ERR_NO_MEMORY         Error = C.SR_ERR_NO_MEMORY
	ERR_NOT_FOUND         Error = C.SR_ERR_NOT_FOUND
	ERR_EXISTS            Error = C.SR_ERR_EXISTS
	ERR_INTERNAL          Error = C.SR_ERR_INTERNAL
	ERR_UNSUPPORTED       Error = C.SR_ERR_UNSUPPORTED
	ERR_VALIDATION_FAILED Error = C.SR_ERR_VALIDATION_FAILED
	ERR_OPERATION_FAILED  Error = C.SR_ERR_OPERATION_FAILED
	ERR_UNAUTHORIZED      Error = C.SR_ERR_UNAUTHORIZED
	ERR_LOCKED            Error = C.SR_ERR_LOCKED
	ERR_TIME_OUT          Error = C.SR_ERR_TIME_OUT
	ERR_CALLBACK_FAILED   Error = C.SR_ERR_CALLBACK_FAILED
	ERR_CALLBACK_SHELVE   Error = C.SR_ERR_CALLBACK_SHELVE
)

var srErrorNames = map[Error]string{
	ERR_OK:                "ERR_OK",
	ERR_INVAL_ARG:         "ERR_INVAL_ARG",
	ERR_LY:                "ERR_LY",
	ERR_SYS:               "ERR_SYS",
	ERR_NO_MEMORY:         "ERR_NO_MEMORY",
	ERR_EXISTS:            "ERR_EXISTS",
	ERR_NOT_FOUND:         "ERR_NOT_FOUND",
	ERR_INTERNAL:          "ERR_INTERNAL",
	ERR_UNSUPPORTED:       "ERR_UNSUPPORTED",
	ERR_VALIDATION_FAILED: "ERR_VALIDATION_FAILED",
	ERR_OPERATION_FAILED:  "ERR_OPERATION_FAILED",
	ERR_UNAUTHORIZED:      "ERR_UNAUTHORIZED",
	ERR_LOCKED:            "ERR_LOCKED",
	ERR_TIME_OUT:          "ERR_TIME_OUT",
	ERR_CALLBACK_FAILED:   "ERR_CALLBACK_FAILED",
	ERR_CALLBACK_SHELVE:   "ERR_CALLBACK_SHELVE",
}

var srErrorValues = map[string]Error{
	"ERR_OK":                ERR_OK,
	"ERR_INVAL_ARG":         ERR_INVAL_ARG,
	"ERR_LY":                ERR_LY,
	"ERR_SYS":               ERR_SYS,
	"ERR_NO_MEMORY":         ERR_NO_MEMORY,
	"ERR_EXISTS":            ERR_EXISTS,
	"ERR_NOT_FOUND":         ERR_NOT_FOUND,
	"ERR_INTERNAL":          ERR_INTERNAL,
	"ERR_UNSUPPORTED":       ERR_UNSUPPORTED,
	"ERR_VALIDATION_FAILED": ERR_VALIDATION_FAILED,
	"ERR_OPERATION_FAILED":  ERR_OPERATION_FAILED,
	"ERR_UNAUTHORIZED":      ERR_UNAUTHORIZED,
	"ERR_LOCKED":            ERR_LOCKED,
	"ERR_TIME_OUT":          ERR_TIME_OUT,
	"ERR_CALLBACK_FAILED":   ERR_CALLBACK_FAILED,
	"ERR_CALLBACK_SHELVE":   ERR_CALLBACK_SHELVE,
}

func (v Error) String() string {
	if s, ok := srErrorNames[v]; ok {
		return s
	}

	return fmt.Sprintf("Error(%d)", v)
}

func ParseError(rc C.int) error {
	if rc == C.SR_ERR_OK {
		return nil
	}

	return fmt.Errorf("ParseError %w : %s", ErrSysrepo, Error(rc))
}
