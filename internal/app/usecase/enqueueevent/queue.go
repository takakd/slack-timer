package enqueueevent

type QueueMessage struct {
	UserId string
}

// Make entities permanent.
type Queue interface {
	Enqueue(message *QueueMessage) (string, error)
}
