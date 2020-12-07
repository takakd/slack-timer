package notifyevent

// Notifier defines notifying methods used in notifying usecase.
type Notifier interface {
	// Notify an event to user.
	Notify(userID string, message string) error
}
