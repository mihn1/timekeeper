package core

import "github.com/mihn1/timekeeper/internal/data"

func SeedDataInMem(t *TimeKeeper) {
	t.categoryStore.AddCategory(data.Category{Id: data.WORK, Name: "Work"})
	t.categoryStore.AddCategory(data.Category{Id: data.ENTERTAINMENT, Name: "Entertainment"})
	t.categoryStore.AddCategory(data.Category{Id: data.UNDEFINED, Name: "Undefined"})
}
