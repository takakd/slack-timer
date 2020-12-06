package enqueueevent

type Queue interface {
	Enqueue(message QueueMessage) (string, error)
}

type QueueMessage struct {
	UserId string
}
