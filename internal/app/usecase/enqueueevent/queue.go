package enqueueevent

// Queue defines enqueueing methods used in enqueue usecase.
type Queue interface {
	Enqueue(message QueueMessage) (string, error)
}

// QueueMessage is parameter of Queue.Enqueue.
type QueueMessage struct {
	UserID string
	Text   string
}
