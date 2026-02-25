package service

import "goboot-cli/goboot/repository"

type Foo interface {
	GetFoo() (int, error)
}

type foo struct {
	repository repository.Foo
}

func NewFoo(
	repository repository.Foo,
) Foo {
	return &foo{
		repository: repository,
	}
}

func (f *foo) GetFoo() (int, error) {
	getFoo, err := f.repository.GetFoo(1)
	if err != nil {
		return 0, err
	}
	return getFoo, nil
}
