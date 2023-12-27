package complaint

import (
	"strings"
)

type State string

const (
	Unknown          State = "Unknown"
	Fined                  = "FINED"
	HearingScheduled       = "Hearing Scheduled"

	NoticeOfDecision = "Notice of Decision"
	ClosedPenalty    = "Plead Guilty"
	// ClosedNotGuilty  = "CLOSED (not guilty)"
	// ClosedGuilty     = "CLOSED (guilty)"
	ClosedInspection = "CLOSED (Referred to S&E)"
	ClosedUnableToID = "CLOSED (Unable to ID)"
	Invalid          = "CLOSED (Invalid)"
	Expired          = "EXPIRED (Unknown)"
)

var AllStates []State = []State{
	Unknown,
	Fined,
	HearingScheduled,
	NoticeOfDecision,
	ClosedPenalty,
	ClosedInspection,
	ClosedUnableToID,
	Invalid,
	Expired,
}

func (s State) String() string {
	return string(s)
}

func DetectState(s string) State {
	switch {
	case strings.Contains(s, "pled guilty") || strings.Contains(s, "STIP violation") || strings.Contains(s, "has paid a penalty") || strings.Contains(s, "pleaded guilty"):
		return ClosedPenalty
	case strings.Contains(s, "Referred to S&E") || strings.Contains(s, "S&E Referral") || strings.Contains(s, "inspected at the TLC inspection facility"):
		return ClosedInspection
	case strings.Contains(s, "scheduled") || strings.Contains(s, "Scheduled") || strings.Contains(s, "hearing sch"):
		return HearingScheduled
	case strings.Contains(s, "mailed to driver") || strings.Contains(s, "sent to driver") || strings.Contains(s, "issued a notice to the driver"):
		return Fined
	case strings.Contains(strings.ToLower(s), "unable to identify") || strings.Contains(s, "not been able to identify the driver"):
		return ClosedUnableToID
	case strings.Contains(s, "Not a TLC violation") || strings.Contains(s, "not a violation of TLC rules") || strings.Contains(s, "not against TLC rules"):
		return Invalid
	case strings.Contains(s, "We have been unable to reach you"):
		return Expired
	}
	return Unknown
}
