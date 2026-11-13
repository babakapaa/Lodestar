package sandbox

import "errors"

var ErrInvalidClass = errors.New("gpu class is invalid")

var ErrClassNotFound = errors.New("gpu class not found")

var ErrClassTaken = errors.New("gpu class name is taken")

var ErrInvalidTerm = errors.New("term is invalid")

var ErrTermNotFound = errors.New("term not found")

var ErrTermOverlap = errors.New("terms overlap")

var ErrNoTerm = errors.New("no term is running now")

var ErrInvalidTransfer = errors.New("transfer is invalid")

var ErrNotInTeam = errors.New("not a member of this team")

var ErrNotEnoughHours = errors.New("not enough gpu hours")

var ErrInvalidTopup = errors.New("top-up request is invalid")

var ErrWorkNotFound = errors.New("work not found")

var ErrTopupNotFound = errors.New("top-up request not found or decided")

var ErrTopupPending = errors.New("top-up request already waits")

var ErrInvalidTemplate = errors.New("environment template is invalid")

var ErrTemplateNotFound = errors.New("environment template not found")

var ErrTemplateTaken = errors.New("environment template name is taken")

var ErrInvalidNode = errors.New("gpu node is invalid")

var ErrNodeNotFound = errors.New("gpu node not found")

var ErrNodeTaken = errors.New("gpu node name is taken")

var ErrNodeBusy = errors.New("gpu node still runs work")

var ErrBadNodeToken = errors.New("gpu node token is wrong")

var ErrInvalidRun = errors.New("run is invalid")

var ErrRunNotFound = errors.New("run not found")

var ErrRunEnded = errors.New("run has already ended")
