package core

type Observer interface {
	Start() error
	Stop() error
}
