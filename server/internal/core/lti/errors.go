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

var ErrNotLinked = errors.New("lti context is not linked")

var ErrSourceFixed = errors.New("assignment source is chosen by the university")

var ErrInvalidSource = errors.New("assignment source is invalid")

var ErrPlatformUnavailable = errors.New("lms answered with an error")

var ErrDeepLinkNotFound = errors.New("deep linking request not found or expired")

var ErrInvalidDeepLink = errors.New("deep linking answer is invalid")
