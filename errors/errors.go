// Package errors has functions and structs to facilitade catch specifics
// errors thougth your program.
//
// For example, use a BadRequest to sinalize a specific error to client
package errors

import "errors"

type Error interface {
	error
	Status() int
}

type BadRequest struct {
	Err error
}

type NotFound struct {
	Err error
}

type InternalServer struct {
	Err error
}

// Allows StatusError to satisfy the error interface.
func (se BadRequest) Error() string {
	return se.Err.Error()
}

func NewBadRequest(err string) BadRequest {
	return BadRequest{Err: errors.New(err)}
}

func (se NotFound) Error() string {
	return se.Err.Error()
}

func (se InternalServer) Error() string {
	return se.Err.Error()
}
