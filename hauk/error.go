package hauk

// SessionExpiredError signals that a given hauk session has expired
type SessionExpiredError struct{}

func (t SessionExpiredError) Error() string { return "Session expired" }

// MalformedResponseError signals that hauk did respond with a malformed response
type MalformedResponseError struct{}

func (t MalformedResponseError) Error() string { return "Malformed response from Hauk" }
