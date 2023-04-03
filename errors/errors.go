// Package errors has functions and structs to facilitade catch specifics
// errors thougth your program.
//
// For example, use a BadRequest to sinalize a specific error to client
package errors

import "errors"

type BadRequest struct {
	Err error
}

type NotFound struct {
	Err error
}

type InternalServerError struct {
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
func (se InternalServerError) Error() string {
	return se.Err.Error()
}

func NewInternalServerError(err string) error {
	return InternalServerError{Err: errors.New(err)}
}

func New(err string) error {
	return errors.New(err)
}
