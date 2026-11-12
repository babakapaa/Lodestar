package audit

import "errors"

var ErrForbidden = errors.New("audit log needs core.audit.view")
