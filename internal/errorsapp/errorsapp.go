package errorsapp

import "errors"

var ErrLinkAlreadyExists = errors.New("the link already exists")

var ErrInvalidLinkReceived = errors.New("an invalid link was received")
