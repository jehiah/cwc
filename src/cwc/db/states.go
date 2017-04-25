package db

import (
	"strings"
)

type State string

const (
	Unknown          State = "Unknown"
	Submitted              = "Submitted"
	Invalid                = "Invalid"
	Fined                  = "FINED"
	ClosedPenalty          = "Plead Guilty"
	HearingScheduled       = "Hearing Scheduled"
	NoticeOfDecision       = "Notice of Decision"
	ClosedGuilty           = "FINED (guilty)"
	ClosedNotGuilty        = "CLOSED (not guilty)"
	ClosedInspection       = "CLOSED (Referred to S&E)"
	ClosedUnableToID       = "CLOSED (Uanble to ID)"
)

func (s State) String() string {
	return string(s)
}

func DetectState(s string) State {
	switch {
	case strings.Contains(s, "pled guilty") || strings.Contains(s, "STIP violation") || strings.Contains(s, "has paid a penalty") || strings.Contains(s, "pleaded guilty"):
		return ClosedPenalty
	case strings.Contains(s, "Referred to S&E"):
		return ClosedInspection
	case strings.Contains(s, "scheduled") || strings.Contains(s, "Scheduled") || strings.Contains(s, "hearing sch"):
		return HearingScheduled
	case strings.Contains(s, "mailed to driver") || strings.Contains(s, "sent to driver"):
		return Fined
	case strings.Contains(s, "unable to identify"):
		return ClosedUnableToID
	case strings.Contains(s, "Not a TLC violation"):
		return Invalid
	}
	return Unknown
}
