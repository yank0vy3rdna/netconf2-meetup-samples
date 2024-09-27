package sysrepo

/*
#cgo LDFLAGS: -lsysrepo -lyang
#include "helper.h"
*/
import "C"

import (
	"unsafe"
)

type Interface interface {
	Connect(opts C.sr_conn_options_t, con **C.sr_conn_ctx_t) error
	Disconnect(con *C.sr_conn_ctx_t) error
	CreateSession(c *Connection) (SessionInterface, error)
	CreateSessionForDatastore(c *Connection, d Datastore) (SessionInterface, error)
	SessionStart(connection *C.sr_conn_ctx_t, datastore C.sr_datastore_t, session **C.sr_session_ctx_t) error
	SessionStop(session *C.sr_session_ctx_t) error
	ModuleChangeSubscribe(
		session *C.sr_session_ctx_t,
		moduleName *C.char,
		xpath *C.char,
		callback C.sr_module_change_cb,
		privateData unsafe.Pointer,
		priority C.uint32_t,
		opts C.uint32_t,
		subscription **C.sr_subscription_ctx_t,
	) error
	Unsubscribe(subscription *C.sr_subscription_ctx_t) error
	GetData(
		session *C.sr_session_ctx_t,
		xpath *C.char,
		maxDepth C.uint32_t,
		timeoutMs C.uint32_t,
		opts C.sr_get_options_t,
	) (Data, error)
	CopyConfig(session *C.sr_session_ctx_t, moduleName *C.char, srcDatastore C.sr_datastore_t, timeoutMs C.uint32_t) error
}

type Data interface {
	DataTreeToString(format LydFormat) string
	Free()
}

type SubscriptionInterface interface {
	Unsubscribe() error
}

type SessionInterface interface {
	Stop() error
	GetDataByModuleName(module string) (Data, error)
	GetDataByXpath(xpath string) (Data, error)
	Subscribe(module string, xpath string, callback ModuleChangeCallback) (SubscriptionInterface, error)
	ID() string
	CloneWithContext(ctx *C.sr_session_ctx_t) SessionInterface
}
