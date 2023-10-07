package storage

import "gorm.io/gorm"

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

func (c *categoryStore) Get(name string) (*Category, error) {
	category := &Category{}
	result := c.db.Where("name = ?", name).First(category)
	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (c *categoryStore) Where(query interface{}, args ...interface{}) ([]Category, error) {
	categories := []Category{}
	result := c.db.Where(query, args).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (c *categoryStore) List() ([]Category, error) {
	categories := []Category{}
	result := c.db.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (c *categoryStore) Update(category Category) error {
	res := c.db.Save(&category)
	return res.Error
}
