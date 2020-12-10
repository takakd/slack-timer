package appcontext

import "time"

// AppContext defines the context interface used in the app.
type AppContext interface {
	RequestID() string
	HandlerCalledTime() time.Time
}
