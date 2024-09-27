package sysrepo

/*
#include "helper.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var ErrCallbackAlreadyRegistered = errors.New("callback for module already exists")

type OriginalSessionCallback struct {
	Session  SessionInterface
	Callback ModuleChangeCallback
}

type ModuleChangeCallback func(session SessionInterface, module string, event NotifyEvent) int

func recoverModuleChangeCallback(cb ModuleChangeCallback) ModuleChangeCallback {
	return func(session SessionInterface, module string, event NotifyEvent) (result int) {
		defer func() {
			if r := recover(); r != nil {
				result = int(ERR_INTERNAL)
			}
		}()

		return cb(session, module, event)
	}
}

type moduleChangeCallbacksStorage struct {
	mu      sync.RWMutex
	storage map[string]map[string]OriginalSessionCallback
}

func (m *moduleChangeCallbacksStorage) Load(module, sessionID string) (OriginalSessionCallback, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	callbacks, ok := m.storage[module]
	if !ok {
		return OriginalSessionCallback{}, false
	}

	callback, ok := callbacks[sessionID]

	return callback, ok
}

func (m *moduleChangeCallbacksStorage) LoadModule(module string) ([]OriginalSessionCallback, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	callbacks, ok := m.storage[module]
	callbacksSlice := make([]OriginalSessionCallback, 0, len(callbacks))
	for _, cb := range callbacks {
		callbacksSlice = append(callbacksSlice, cb)
	}

	return callbacksSlice, ok
}

func (m *moduleChangeCallbacksStorage) Store(module string, session SessionInterface, callback ModuleChangeCallback) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	callbacks, ok := m.storage[module]
	if !ok {
		newCallbacks := make(map[string]OriginalSessionCallback)
		m.storage[module] = newCallbacks
		callbacks = newCallbacks
	}

	if _, exists := callbacks[session.ID()]; exists {
		return fmt.Errorf("%w: %s, sessionID %s", ErrCallbackAlreadyRegistered, module, session.ID())
	}

	sessionCallback := OriginalSessionCallback{
		Session:  session,
		Callback: recoverModuleChangeCallback(callback),
	}
	callbacks[session.ID()] = sessionCallback

	return nil
}

func (m *moduleChangeCallbacksStorage) Delete(module, sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	callbacks, ok := m.storage[module]
	if !ok {
		return
	}

	delete(callbacks, sessionID)
	if len(callbacks) == 0 {
		delete(m.storage, module)
	}
}

func (m *moduleChangeCallbacksStorage) DeleteModule(module string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.storage, module)
}

var moduleChangeCallbacks = moduleChangeCallbacksStorage{
	storage: make(map[string]map[string]OriginalSessionCallback),
	mu:      sync.RWMutex{},
}

func ModuleChangeCallbackRegister(module string, session SessionInterface, callback ModuleChangeCallback) error {
	return moduleChangeCallbacks.Store(module, session, callback)
}

func ModuleChangeCallbackUnregister(module string, sessionIDs ...string) {
	if len(sessionIDs) == 0 {
		moduleChangeCallbacks.DeleteModule(module)
	}

	for _, sessionID := range sessionIDs {
		moduleChangeCallbacks.Delete(module, sessionID)
	}
}

func ModuleChangeCallbackGet(module string) ([]OriginalSessionCallback, bool) {
	return moduleChangeCallbacks.LoadModule(module)
}

//export go_sr_module_change_cb
func go_sr_module_change_cb(
	cSess *C.sr_session_ctx_t,
	subID uint32,
	cModule *C.char,
	cXPath *C.char,
	cEvent C.sr_event_t,
	requestID uint32,
	cData unsafe.Pointer,
) int {
	module := C.GoString(cModule)
	result := int(ERR_OK)

	if callbacks, ok := ModuleChangeCallbackGet(module); ok {
		for _, callback := range callbacks {
			session := callback.Session.CloneWithContext(cSess)

			rVal := callback.Callback(session, module, NotifyEvent(cEvent))

			if rVal != int(ERR_OK) {
				result = rVal
			}
		}

		return result
	}

	return int(ERR_INTERNAL)
}
