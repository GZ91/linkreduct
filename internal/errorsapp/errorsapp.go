package errorsapp

import "errors"

var ErrLinkAlreadyExists = errors.New("the link already exists")

var ErrInvalidLinkReceived = errors.New("an invalid link was received")

var ErrLineURLDeleted = errors.New("this entry has been deleted")
