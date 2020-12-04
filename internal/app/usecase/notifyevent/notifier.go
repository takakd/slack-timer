package notifyevent

type Notifier interface {
	// Notify an event to user.
	Notify(userId string, message string) error
}
