//nolint:nlreturn
package sysrepo

/*
#include "helper.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	ErrNodeNotFound = errors.New("failed to find YANG node")

	errWrongParams = errors.New("too many params provided as datastore")
)

type Wrapper struct{}

func (w *Wrapper) Connect(opts C.sr_conn_options_t, con **C.sr_conn_ctx_t) error {
	rc := C.sr_connect(opts, con)

	return ParseError(rc)
}

func (w *Wrapper) Disconnect(con *C.sr_conn_ctx_t) error {
	rc := C.sr_disconnect(con)

	return ParseError(rc)
}

func (w *Wrapper) CreateSession(c *Connection) (SessionInterface, error) {
	return NewSession(c, DS_RUNNING)
}

func (w *Wrapper) CreateSessionForDatastore(c *Connection, d Datastore) (SessionInterface, error) {
	return NewSession(c, d)
}

func (w *Wrapper) SessionStart(connection *C.sr_conn_ctx_t, datastore C.sr_datastore_t, session **C.sr_session_ctx_t) error {
	rc := C.sr_session_start(connection, datastore, session)

	return ParseError(rc)
}

func (w *Wrapper) SessionStop(session *C.sr_session_ctx_t) error {
	rc := C.sr_session_stop(session)

	return ParseError(rc)
}

func (w *Wrapper) ModuleChangeSubscribe(
	session *C.sr_session_ctx_t,
	moduleName *C.char,
	xpath *C.char,
	callback C.sr_module_change_cb,
	privateData unsafe.Pointer,
	priority C.uint32_t,
	opts C.uint32_t,
	subscription **C.sr_subscription_ctx_t,
) error {
	rc := C.sr_module_change_subscribe(session, moduleName, xpath, callback, privateData, priority, opts, subscription)

	return ParseError(rc)
}

func (w *Wrapper) Unsubscribe(subscription *C.sr_subscription_ctx_t) error {
	rc := C.sr_unsubscribe(subscription)

	return ParseError(rc)
}

func (w *Wrapper) GetData(
	session *C.sr_session_ctx_t,
	xpath *C.char,
	maxDepth C.uint32_t,
	timeoutMs C.uint32_t,
	opts C.sr_get_options_t,
) (Data, error) {
	var d *C.sr_data_t
	rc := C.sr_get_data(session, xpath, maxDepth, timeoutMs, opts, &d)

	if err := ParseError(rc); err != nil {
		return nil, err
	}

	if d == nil {
		return nil, ErrNodeNotFound
	}

	return (*data)(d), nil
}

func (w *Wrapper) CopyConfig(
	session *C.sr_session_ctx_t,
	moduleName *C.char,
	srcDatastore C.sr_datastore_t,
	timeoutMs C.uint32_t,
) error {
	rc := C.sr_copy_config(session, moduleName, srcDatastore, timeoutMs)

	return ParseError(rc)
}
