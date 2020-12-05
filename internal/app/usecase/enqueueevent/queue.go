package enqueueevent

// Make entities permanent.
type Queue interface {
	Enqueue(message *QueueMessage) (string, error)
}

type QueueMessage struct {
	UserId string
}
