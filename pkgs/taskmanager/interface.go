package taskmanager

////////////////////////////////////////////////////////////////////////////////

//go:generate mockgen -source=interface.go -destination=taskmanager_mock.go -package=taskmanager
type Task interface {
	Execute() bool
	SetRetrySignal() <-chan struct{}
	GetID() string
	GetProgress() int64
	Cancel()
}
