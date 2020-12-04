package enqueueevent

type OutputData struct {
	Result             error
	NotifiedUserIdList []string
	QueueMessageIdList []string
}

type OutputPort interface {
	Output(data OutputData)
}
