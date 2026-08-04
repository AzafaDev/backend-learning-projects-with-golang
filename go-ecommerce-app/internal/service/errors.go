package service

import "errors"

var ErrEmailTaken = errors.New("email is already registered")
var ErrExpiredRefreshToken = errors.New("expired refresh token")
