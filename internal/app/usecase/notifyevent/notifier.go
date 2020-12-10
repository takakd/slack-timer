package notifyevent

import "slacktimer/internal/app/util/appcontext"

// Notifier defines notifying methods used in notifying usecase.
type Notifier interface {
	// Notify an event to user.
	Notify(ac appcontext.AppContext, userID string, message string) error
}
