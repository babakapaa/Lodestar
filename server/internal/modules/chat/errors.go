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
