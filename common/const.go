package common

import "time"

const (
	EVENT_MESSAGE        = "message"
	EVENT_EXECUTE        = "execute"
	EVENT_JOIN           = "join"
	EVENT_LEAVE          = "leave"
	EVENT_LEAVE_ALL      = "leave_all"
	GENERAL_ROOM         = "general"
	DEFAULT_IGNORE_ID    = ""
	CURRENT_USER         = "currentUser"
	DEFAULT_TIMEOUT_SSH  = time.Second * 30
	DEFAULT_TIMEOUT_DIAL = time.Second * 2
)
