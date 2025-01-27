package core

type Observer interface {
	StartObserving(t *TimeKeeper) error
}
