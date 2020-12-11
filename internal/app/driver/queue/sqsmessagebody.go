package queue

// SqsMessageBody is message body structure enqueued in AWS SQS.
type SqsMessageBody struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

// NewSqsMessageBody creates new struct.
func NewSqsMessageBody() *SqsMessageBody {
	return &SqsMessageBody{}
}

// NOTE: Redundancy
//// JSON returns JSON string.
//func (s SqsMessageBody)JSON() (js string, err error){
//	b, err := json.Marshal(s)
//	if err != nil {
//		return
//	}
//	js = string(b)
//	return
//}
