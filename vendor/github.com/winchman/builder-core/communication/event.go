package comm

// EventType is a type for constants that indicate the type of event reported
type EventType uint8

// Event is the type that will be sent over the event channel
type Event interface {
	EventType() EventType
	Data() map[string]interface{}
}

const (
	// RequestedEvent is for when a build is initially requested
	RequestedEvent EventType = iota

	// BuildEvent is for when a `docker build` command starts
	BuildEvent

	// BuildCompletedEvent is for when a `docker build` command completes
	BuildCompletedEvent

	// TagEvent is for when a `docker tag` command starts
	TagEvent

	// TagCompletedEvent is for when a `docker tag` command finishes
	TagCompletedEvent

	// PushEvent is for when a `docker push` command starts
	PushEvent

	// PushCompletedEvent is for when a `docker push` command finishes
	PushCompletedEvent

	// CompletedEvent is for whe nthe entire build finishes (corresopnds to a RequestedEvent)
	CompletedEvent
)

func (t EventType) String() string {
	switch t {
	case RequestedEvent:
		return "RequestedEvent"
	case CompletedEvent:
		return "CompletedEvent"
	case BuildEvent:
		return "BuildEvent"
	case BuildCompletedEvent:
		return "BuildCompletedEvent"
	case TagEvent:
		return "TagEvent"
	case TagCompletedEvent:
		return "TagCompletedEvent"
	case PushEvent:
		return "PushEvent"
	case PushCompletedEvent:
		return "PushCompletedEvent"
	default:
		return ""
	}
}

type event struct {
	eventType EventType
	data      map[string]interface{}
}

func (e *event) EventType() EventType         { return e.eventType }
func (e *event) Data() map[string]interface{} { return e.data }
