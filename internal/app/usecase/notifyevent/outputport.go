package notifyevent

import "slacktimer/internal/app/util/appcontext"

// OutputPort defines outputport method of notifying events usecase.
type OutputPort interface {
	Output(ac appcontext.AppContext, data OutputData)
}

// OutputData is parameter of OutputPort.
type OutputData struct {
	// nil if success, otherwise error.
	Result error
	// Notified UserID.
	UserID string
}
