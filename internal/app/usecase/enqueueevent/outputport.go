package enqueueevent

// OutputPort defines outputport method of enqueueing events usecase.
type OutputPort interface {
	Output(data OutputData)
}

// OutputData is parameter of OutputPort.
type OutputData struct {
	Result             error
	// Succeeded in Notifying ID list.
	NotifiedUserIDList []string
	QueueMessageIDList []string
}
