package errorsapp

import "errors"

// ErrLinkAlreadyExists представляет ошибку, указывающую на то, что ссылка уже существует.
var ErrLinkAlreadyExists = errors.New("the link already exists")

// ErrInvalidLinkReceived представляет ошибку, указывающую на то, что получена некорректная ссылка.
var ErrInvalidLinkReceived = errors.New("an invalid link was received")

// ErrLineURLDeleted представляет ошибку, указывающую на то, что данная запись была удалена.
var ErrLineURLDeleted = errors.New("this entry has been deleted")
