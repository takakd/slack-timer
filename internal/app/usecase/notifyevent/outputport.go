package notifyevent

type OutputPort interface {
	Output(data OutputData)
}

type OutputData struct {
	// nil if success, otherwise error.
	Result error
	// Notified UserID.
	UserId string
}
