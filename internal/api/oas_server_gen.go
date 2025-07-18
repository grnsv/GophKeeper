// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// LoginPost implements POST /login operation.
	//
	// Authenticate user.
	//
	// POST /login
	LoginPost(ctx context.Context, req *UserCredentials) (LoginPostRes, error)
	// RecordsGet implements GET /records operation.
	//
	// Get all user records.
	//
	// GET /records
	RecordsGet(ctx context.Context) (RecordsGetRes, error)
	// RecordsIDDelete implements DELETE /records/{id} operation.
	//
	// Delete record.
	//
	// DELETE /records/{id}
	RecordsIDDelete(ctx context.Context, params RecordsIDDeleteParams) (RecordsIDDeleteRes, error)
	// RecordsIDGet implements GET /records/{id} operation.
	//
	// Get specific record.
	//
	// GET /records/{id}
	RecordsIDGet(ctx context.Context, params RecordsIDGetParams) (RecordsIDGetRes, error)
	// RecordsIDPut implements PUT /records/{id} operation.
	//
	// Create or update record.
	//
	// PUT /records/{id}
	RecordsIDPut(ctx context.Context, req *Record, params RecordsIDPutParams) (RecordsIDPutRes, error)
	// RegisterPost implements POST /register operation.
	//
	// Register new user.
	//
	// POST /register
	RegisterPost(ctx context.Context, req *UserCredentials) (RegisterPostRes, error)
	// VersionGet implements GET /version operation.
	//
	// Get server version.
	//
	// GET /version
	VersionGet(ctx context.Context) (*VersionInfo, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
