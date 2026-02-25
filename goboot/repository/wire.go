package repository

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewFoo,
)
