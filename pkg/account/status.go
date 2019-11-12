package account

// Status is an account status type
type Status string

const (
	// None status
	None Status = "None"
	// Ready status
	Ready Status = "Ready"
	// NotReady status
	NotReady Status = "NotReady"
	// Leased status
	Leased Status = "Leased"
	// Orphaned status
	Orphaned Status = "Orphaned"
)
