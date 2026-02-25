package controller

import (
	"goboot-cli/goboot/service"

	"github.com/gin-gonic/gin"
)

type Foo struct {
	fooService service.Foo
}

func (ctl *Foo) GetFoo(c *gin.Context) {
	foo, err := ctl.fooService.GetFoo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": foo,
	})
}
