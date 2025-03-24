package models

type CategoryId int

const (
	EXCLUDED      CategoryId = 0
	WORK          CategoryId = 1
	ENTERTAINMENT CategoryId = 2
	PERSONAL      CategoryId = 3
	UNDEFINED     CategoryId = 4
)

type Category struct {
	Id          CategoryId
	Name        string
	Description string
}
