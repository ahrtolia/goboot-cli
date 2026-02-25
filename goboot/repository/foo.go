package repository

import "gorm.io/gorm"

type Foo interface {
	GetFoo(id int) (int, error)
}
type foo struct {
	db *gorm.DB
}

func NewFoo(
	db *gorm.DB,
) Foo {
	return &foo{
		db: db,
	}
}

func (f *foo) GetFoo(id int) (int, error) {
	return id, nil
}
