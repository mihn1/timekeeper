package models

type CategoryId string

const (
	WORK          CategoryId = "work"
	ENTERTAINMENT CategoryId = "entertainment"
	PERSONAL      CategoryId = "personal"
	UNDEFINED     CategoryId = "undefined"
)

type Category struct {
	Id             CategoryId
	Name           string
	Description    string
	CategoryTypeId int
}
