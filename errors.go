package smszenziva

import "errors"

// List of errors in the Zenziva package.
var (
	ErrNilArgs          = errors.New("nil argument")
	ErrFailedToSendSMS  = errors.New("failed to send SMS")
	ErrEmptyUserKey     = errors.New("empty user key")
	ErrEmptyPasswordKey = errors.New("empty password key")
)
