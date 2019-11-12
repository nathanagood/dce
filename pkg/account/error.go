package account

// StatusTransitionError means that we failed to transition
// an Account or Lease from one status to another,
// likely because the prevStatus condition was not met
type StatusTransitionError struct {
	err string
}

func (e *StatusTransitionError) Error() string {
	return e.err
}

// LeasedError is returned when a consumer attempts to delete an account that is currently at status Leased
type LeasedError struct {
	err string
}

func (e *LeasedError) Error() string {
	return e.err
}
