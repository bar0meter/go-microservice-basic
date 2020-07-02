package notifications

type Dispatcher interface {
	Dispatch() (bool, error)
}
