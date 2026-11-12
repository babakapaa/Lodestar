package access

import "errors"

var ErrForbidden = errors.New("permission denied")

var ErrHierarchy = errors.New("role is not below your highest role")

var ErrPermissionNotHeld = errors.New("you cannot grant a permission you do not have")

var ErrUnknownPermission = errors.New("unknown permission")

var ErrBadPermission = errors.New("permission key is malformed or duplicated")

var ErrRoleNotFound = errors.New("role not found")

var ErrAssignmentNotFound = errors.New("role assignment not found")

var ErrUserNotFound = errors.New("user not found")

var ErrInvalidRole = errors.New("role name, color or icon is invalid")

var ErrInvalidScope = errors.New("scope is invalid")

var ErrInvalidPosition = errors.New("position is out of range")

var ErrDefaultRole = errors.New("default role cannot be moved, deleted or assigned")
