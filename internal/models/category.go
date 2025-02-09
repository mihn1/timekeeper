package models

const (
	WORK          CategoryId = "work"
	ENTERTAINMENT            = "entertainment"
	PERSONAL                 = "personal"
	UNDEFINED                = "undefined"
)

type CategoryId string

type Category struct {
	Id             CategoryId
	Name           string
	CategoryTypeId int
}
