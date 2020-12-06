package enqueueevent

type OutputPort interface {
	Output(data OutputData)
}

type OutputData struct {
	Result             error
	NotifiedUserIdList []string
	QueueMessageIdList []string
}
