package notifyevent

import "slacktimer/internal/app/util/appcontext"

// Notifier defines notifying methods used in notifying usecase.
type Notifier interface {
	// SendMessage an event to user.
	SendMessage(ac appcontext.AppContext, userID string, text string) error
}
