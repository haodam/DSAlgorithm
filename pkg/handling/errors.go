package handling

import "errors"

var (
	ErrBadRequest      = errors.New("bad Request")
	ErrInternalFailure = errors.New("internal failure")
	ErrNotfound        = errors.New("not found")
)

type Error struct {
	appErr   error
	svcError error
}

func (e Error) Error() string {
	return errors.Join(e.svcError, e.appErr).Error()
}

func NewError(svcErr, appErr errorr) error {
	return Error{
		svcError: svcErr,
		appErr:   appErr,
	}
}
