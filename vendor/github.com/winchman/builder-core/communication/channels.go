package comm

type (
	// LogChan is a channel for log entries
	LogChan chan LogEntry

	// EventChan is a channel for status updates
	EventChan chan Event

	// ExitChan is a channel for receiving the final exit value (error or nil)
	ExitChan chan error
)
