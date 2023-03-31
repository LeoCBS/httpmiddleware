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

type ServerError struct {
	Err error
}

// Satisfy the error interface.
func (se BadRequest) Error() string {
	return se.Err.Error()
}

func NewBadRequest(err string) error {
	return BadRequest{Err: errors.New(err)}
}

// Satisfy the error interface.
func (se NotFound) Error() string {
	return se.Err.Error()
}

func NewNotFound(err string) error {
	return NotFound{Err: errors.New(err)}
}

// Satisfy the error interface.
func (se ServerError) Error() string {
	return se.Err.Error()
}

func NewServerError(err string) error {
	return ServerError{Err: errors.New(err)}
}
