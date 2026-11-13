package chat

import "errors"

var ErrChatNotFound = errors.New("chat not found")

var ErrInvalidChat = errors.New("chat kind, title or members are invalid")

var ErrReadOnly = errors.New("chat is archived or only admins write here")

var ErrInvalidMessage = errors.New("message is empty, too long or replies to a missing message")

var ErrMessageNotFound = errors.New("message not found")

var ErrOriginNotFound = errors.New("chat origin not found")

var ErrSyncedMember = errors.New("member comes from the course list")

var ErrOwnerLeaves = errors.New("owner cannot leave the chat")

var ErrBadFile = errors.New("file is missing, not uploaded or not yours")

var ErrBlocked = errors.New("this person does not accept your messages")

var ErrInvalidReport = errors.New("report reason or messages are invalid")

var ErrReportNotFound = errors.New("report not found")

var ErrCallsUnavailable = errors.New("calls are not configured on this server")

var ErrCallsOff = errors.New("calls are turned off by the university")

var ErrCallNotFound = errors.New("call not found")

var ErrCallFull = errors.New("call already has the most people allowed")
