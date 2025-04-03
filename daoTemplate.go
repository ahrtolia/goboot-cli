package main

const daoTemplate = `
package dao

type UserDao struct {}

// FindUserById 根据 ID 查找用户
func (dao *UserDao) FindUserById(id int) string {
    return "User from DB"
}
`
