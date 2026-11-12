package auth

import "errors"

var ErrInvalidToken = errors.New("token is invalid or expired")

var ErrTokenReused = errors.New("refresh token was already used")

var ErrProviderNotFound = errors.New("identity provider not found")

var ErrInvalidProvider = errors.New("identity provider settings are invalid")

var ErrInvalidState = errors.New("login attempt is unknown or expired")

var ErrRedirectNotAllowed = errors.New("redirect uri is not allowed")

var ErrInvalidChallenge = errors.New("code challenge must be S256")

var ErrSessionNotFound = errors.New("session not found")

var ErrLoginFailed = errors.New("identity provider rejected the login")
