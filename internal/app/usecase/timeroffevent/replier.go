package timeroffevent

import "slacktimer/internal/app/util/appcontext"

// Replier defines replying methods used in updatetimerevent usecase.
type Replier interface {
	// SendMessage an event to user.
	SendMessage(ac appcontext.AppContext, userID string, text string) error
}
