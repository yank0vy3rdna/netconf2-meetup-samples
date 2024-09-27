package sysrepo

/*
#include "helper.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	uuid "github.com/google/uuid"
)

const (
	operationTimeout uint32 = 5000  // milliseconds
	commitTimeout    uint32 = 60000 // milliseconds
)

type Datastore C.sr_datastore_t

func (v Datastore) C() C.sr_datastore_t {
	return (C.sr_datastore_t)(v)
}

const (
	DS_STARTUP     Datastore = C.SR_DS_STARTUP
	DS_RUNNING     Datastore = C.SR_DS_RUNNING
	DS_CANDIDATE   Datastore = C.SR_DS_CANDIDATE
	DS_OPERATIONAL Datastore = C.SR_DS_OPERATIONAL
)

type Session struct {
	sysrepoIf  Interface
	connection *Connection
	context    *C.sr_session_ctx_t
	id         string
}

type subscription struct {
	sysrepoIf Interface
	session   SessionInterface
	context   *C.sr_subscription_ctx_t
	module    string
}

func NewSession(connection *Connection, datastore Datastore) (*Session, error) {
	var ctx *C.sr_session_ctx_t
	if err := connection.sysrepoIf.SessionStart(connection.context, datastore.C(), &ctx); err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	id := uuid.New().String()

	res := &Session{
		sysrepoIf:  connection.sysrepoIf,
		connection: connection,
		id:         id,
		context:    ctx,
	}

	return res, nil
}

func (s *Session) Stop() error {
	return s.sysrepoIf.SessionStop(s.context)
}

func (s *Session) GetDataByModuleName(module string) (Data, error) {
	return s.GetDataByXpath(fmt.Sprintf("/%s:*", module))
}

func (s *Session) GetDataByXpath(xpath string) (Data, error) {
	cXPath := C.CString(xpath)
	defer C.free(unsafe.Pointer(cXPath))

	unlimitedDepth := C.uint32_t(0)
	timeoutMs := C.uint32_t(operationTimeout)
	options := C.uint32_t(0)

	data, err := s.sysrepoIf.GetData(s.context, cXPath, unlimitedDepth, timeoutMs, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get data in session %v at %v: %w", s.id, xpath, err)
	}

	return data, nil
}

func (s *Session) Subscribe(module string, xpath string, callback ModuleChangeCallback) (SubscriptionInterface, error) {
	if err := ModuleChangeCallbackRegister(module, s, callback); err != nil {

		return nil, err
	}

	cModule := C.CString(module)
	defer C.free(unsafe.Pointer(cModule))

	cXpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(cXpath))

	priority := C.uint32_t(0)

	var ctx *C.sr_subscription_ctx_t

	err := s.sysrepoIf.ModuleChangeSubscribe(s.context, cModule, cXpath, C.sr_module_change_cb(C._sr_module_change_cb),
		nil, priority, C.uint(SUBSCR_DEFAULT), &ctx)
	if err != nil {
		ModuleChangeCallbackUnregister(module)

		return nil, err
	}

	res := &subscription{
		sysrepoIf: s.sysrepoIf,
		session:   s,
		context:   ctx,
		module:    module,
	}
	res.module = module

	return res, nil
}

func (s *subscription) Unsubscribe() error {
	ModuleChangeCallbackUnregister(s.module)
	if err := s.sysrepoIf.Unsubscribe(s.context); err != nil {
		return err
	}

	return nil
}

func (s *Session) CopyConfig(datastore Datastore) error {
	timeoutMs := C.uint32_t(commitTimeout)

	if err := s.sysrepoIf.CopyConfig(s.context, nil, datastore.C(), timeoutMs); err != nil {
		return err
	}

	return nil
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) CloneWithContext(ctx *C.sr_session_ctx_t) SessionInterface {
	return &Session{
		sysrepoIf:  s.sysrepoIf,
		connection: s.connection,
		context:    ctx,
		id:         s.id,
	}
}
