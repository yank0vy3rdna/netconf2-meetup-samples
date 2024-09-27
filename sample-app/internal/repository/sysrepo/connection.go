package sysrepo

/*
#include "helper.h"
*/
import "C"

import (
	"fmt"
)

// Connection represents C++ struct C.sr_conn_ctx_t in Go
type Connection struct {
	sysrepoIf Interface
	context   *C.sr_conn_ctx_t
}

// NewConnection creates new instance of Connection
func NewConnection(sysrepoIf Interface) (*Connection, error) {

	var ctx *C.sr_conn_ctx_t
	if err := sysrepoIf.Connect(C.uint(CONN_DEFAULT.C()), &ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to sysrepo: %w", err)
	}

	res := &Connection{
		sysrepoIf: sysrepoIf,
		context:   ctx,
	}

	return res, nil
}

// Disconnect disconnects connection
func (c *Connection) Disconnect() error {
	return c.sysrepoIf.Disconnect(c.context)
}

// StartSession starts new session on specified connection
func (c *Connection) StartSession() (SessionInterface, error) {
	return c.sysrepoIf.CreateSession(c)
}

// StartSessionForDatastore starts new session on specified connection with specified Datastore
func (c *Connection) StartSessionForDatastore(d Datastore) (SessionInterface, error) {
	return c.sysrepoIf.CreateSessionForDatastore(c, d)
}
