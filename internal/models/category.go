package models

const (
	WORK CategoryId = iota
	ENTERTAINMENT
	UNDEFINED
) 

type CategoryId int

type Category struct {
	Id             CategoryId
	Name           string
	CategoryTypeId int
}
