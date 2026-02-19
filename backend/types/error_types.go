package types

import (
	"fmt"

	"github.com/google/uuid"
)

/*
// Quiz Service erors
*/
type ErrSessionNotFound struct {
	SessionID uuid.UUID
}

func (e *ErrSessionNotFound) Error() string {
	if e.SessionID == uuid.Nil {
		return fmt.Sprintf("Cannot find session: %v", e.SessionID.String())
	}
	return "Cannot find session"
}

type ErrSessionExpired struct {
	SessionID uuid.UUID
}

func (e *ErrSessionExpired) Error() string {
	if e.SessionID == uuid.Nil {
		return fmt.Sprintf("Session is expired: %v", e.SessionID.String())
	}
	return "Session is expired"
}

type ErrProblemNotFound struct {
	ProblemId uuid.UUID
}

func (e *ErrProblemNotFound) Error() string {
	if e.ProblemId == uuid.Nil {
		return fmt.Sprintf("Cannot find problem: %v", e.ProblemId.String())
	}
	return "Cannot find problem"
}

/*
// End Quiz Service errors
*/
