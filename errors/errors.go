// Package errors has functions and structs to facilitade catch specifics
// errors thougth your program.
//
// For example, use BadRequestError to sinalize a specific error to client
package errors

type Error interface {
	error
	Status() int
}

type BadRequestError struct {
	Err error
}

type NotFoundError struct {
	Err error
}

type InternalServerError struct {
	Err error
}

// Allows StatusError to satisfy the error interface.
func (se BadRequestError) Error() string {
	return se.Err.Error()
}

func (se NotFoundError) Error() string {
	return se.Err.Error()
}

func (se InternalServerError) Error() string {
	return se.Err.Error()
}
