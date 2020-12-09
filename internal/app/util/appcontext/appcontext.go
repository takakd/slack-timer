package appcontext

// AppContext defines the context interface used in the app.
type AppContext interface {
	RequestID() string
}
