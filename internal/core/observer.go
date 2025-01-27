package core

type Observer interface {
	StartObserving(ch chan AppSwitchEvent) error
}
