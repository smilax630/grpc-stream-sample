package gateway

import "github.com/grpc-streamer/pkg/entity"

//ErrNotFoundSimplyWallet represent not found err.
type ErrNotFoundSimplyWallet struct {
	err error
}

//Error method is full of error interface.
func (e *ErrNotFoundSimplyWallet) Error() string {
	return e.err.Error()
}

//ErrNotFoundSimplyFacade represent not found err.
type ErrNotFoundSimplyFacade struct {
	err error
}

//Error method is full of error interface.
func (e *ErrNotFoundSimplyFacade) Error() string {
	return e.err.Error()
}

//ErrNotFound represent not found err.
type ErrNotFound struct {
	err error
}

//Error method is full of error interface.
func (e *ErrNotFound) Error() string {
	return e.err.Error()
}

//ErrAlreadyExistSimplyFacade represent validation err.
type ErrAlreadyExistSimplyFacade struct {
	err error
}

//Error method is full of error interface.
func (e *ErrAlreadyExistSimplyFacade) Error() string {
	return e.err.Error()
}

//ErrBadRequest represent canceled err.
type ErrBadRequest struct {
	err error
}

//Error method is full of error interface.
func (e *ErrBadRequest) Error() string {
	return e.err.Error()
}

//ErrInvalidArgument represent invalid tx err.
type ErrInvalidArgument struct {
	err error
}

//Error method is full of error interface.
func (e *ErrInvalidArgument) Error() string {
	return e.err.Error()
}

//ErrInvalidImageSize represent invalid tx err.
type ErrInvalidImageSize struct {
	err error
}

//toErrInvalidWallet represent invalid tx err.
type ErrInvalidReduceCoin struct {
	err error
}

//Error method is full of error interface.
func (e *ErrInvalidReduceCoin) Error() string {
	return e.err.Error()
}

//Error method is full of error interface.
func (e *ErrInvalidImageSize) Error() string {
	return e.err.Error()
}

func toGatewayError(err error) error {
	entity.Logger.Error(err.Error())

	return err
}
