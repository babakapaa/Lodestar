package lti

import "errors"

var ErrPlatformNotFound = errors.New("lti platform not found")

var ErrInvalidPlatform = errors.New("lti platform fields are invalid")

var ErrPlatformTaken = errors.New("lti platform with this issuer and client id already exists")

var ErrPlatformAmbiguous = errors.New("several lti platforms share this issuer")

var ErrContextNotFound = errors.New("lti context not found")

var ErrUnknownLinkKind = errors.New("unknown link kind")

var ErrInvalidState = errors.New("lti launch state is invalid or expired")

var ErrLaunchRejected = errors.New("lti launch rejected")
