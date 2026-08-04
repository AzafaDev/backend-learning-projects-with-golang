package service

import "errors"

var ErrEmailTaken = errors.New("email is already registered")
