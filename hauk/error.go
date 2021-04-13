package hauk

// SessionExpiredError signals that a given hauk session has expired
type SessionExpiredError struct{}

func (t *SessionExpiredError) Error() string { return "Session expired" }
