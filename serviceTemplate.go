package main

const serviceTemplate = `
package service

type UserService struct {}

// GetUser 获取用户信息
func (s *UserService) GetUser() string {
    return "John Doe"
}

// CreateUser 创建一个新的用户
func (s *UserService) CreateUser() {
    // 这里可以是实际的逻辑，例如保存到数据库
}
`
