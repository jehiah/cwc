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
	ClosedGuilty           = "FINED (guilty)"
	ClosedNotGuilty        = "CLOSED (not guilty)"
	ClosedInspection       = "CLOSED (Referred to S&E)"
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
	case strings.Contains(s, "scheduled"):
		return HearingScheduled
	case strings.Contains(s, "mailed to driver"):
		return Fined
	}
	return Unknown
}
