package storage

import "gorm.io/gorm"

//nolint:gochecknoglobals
var categoryList = []string{
	"Groceries",
	"Dining",
	"Shopping",
	"Transport",
	"Fuel",
	"Other",
	"Income",
}

type categoryStore struct {
	db *gorm.DB
}

var _ CategoryGetter = &categoryStore{}

func (c *categoryStore) Get(category *Category) error {
	result := c.db.Where(&category).First(category)
	return result.Error
}

func (c *categoryStore) List() ([]*Category, error) {
	categories := []*Category{}

	result := c.db.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (c *categoryStore) Update(category *Category) error {
	res := c.db.Save(&category)

	return res.Error
}
