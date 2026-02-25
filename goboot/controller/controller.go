package controller

import (
	"github.com/ahrtolia/goboot/pkg/gin_starter"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	ginServer *gin_starter.Server,
	fooController Foo,
) (*gin.Engine, error) {

	r := ginServer.GetRouter()

	r.GET("/foo", fooController.GetFoo)

	return r, nil
}
