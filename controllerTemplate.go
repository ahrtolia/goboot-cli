package main

const controllerTemplate = `
package controller

import (
    "github.com/gin-gonic/gin"
	"%s/pkg/service"
)

type UserController struct {
    UserService service.UserService
}

// GetUser 用于获取用户信息的接口
func (ctrl *UserController) GetUser(c *gin.Context) {
    // 模拟业务逻辑处理
    user := ctrl.UserService.GetUser()
    c.JSON(200, user)
}

// CreateUser 用于创建用户的接口
func (ctrl *UserController) CreateUser(c *gin.Context) {
    // 模拟业务逻辑处理
    ctrl.UserService.CreateUser()
    c.JSON(201, gin.H{"message": "User created"})
}
`
