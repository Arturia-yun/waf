package models

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

// User 用户模型
type User struct {
	Id         int64
	Username   string    `xorm:"varchar(50) notnull unique"`
	Password   string    `xorm:"varchar(32) notnull"`
	LastLogin  time.Time `xorm:"DateTime"`
	LastChange time.Time `xorm:"updated"`
	Version    int       `xorm:"version"` // 乐观锁
}

// ListUser 获取所有用户
func ListUser() (users []User, err error) {
	users = make([]User, 0)
	err = Engine.Find(&users)
	return users, err
}

// GetUserById 根据ID获取用户
func GetUserById(id int64) (*User, error) {
	user := &User{Id: id}
	has, err := Engine.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

// GetUserByName 根据用户名获取用户
func GetUserByName(username string) (*User, error) {
	user := &User{Username: username}
	has, err := Engine.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

// NewUser 创建新用户
func NewUser(username, password string) (err error) {
	// 对密码进行MD5加密
	md5Pwd := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	
	user := &User{
		Username: username,
		Password: md5Pwd,
	}
	
	_, err = Engine.Insert(user)
	return err
}

// UpdateUser 更新用户信息
func UpdateUser(id int64, username, password string) (err error) {
	user := new(User)
	has, err := Engine.Id(id).Get(user)
	if err != nil {
		return err
	}
	
	if !has {
		return fmt.Errorf("用户不存在")
	}
	
	user.Username = username
	
	// 如果提供了新密码，则更新密码
	if password != "" {
		md5Pwd := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		user.Password = md5Pwd
	}
	
	_, err = Engine.Id(id).Update(user)
	return err
}

// DelUser 删除用户
func DelUser(id int64) (err error) {
	_, err = Engine.Delete(&User{Id: id})
	return err
}

// ValidateUser 验证用户登录
func ValidateUser(username, password string) (user *User, ok bool) {
	md5Pwd := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	
	user = &User{
		Username: username,
		Password: md5Pwd,
	}
	
	has, err := Engine.Get(user)
	if err != nil {
		log.Printf("验证用户失败: %v", err)
		return nil, false
	}
	
	if has {
		// 更新最后登录时间
		user.LastLogin = time.Now()
		_, err = Engine.Id(user.Id).Cols("last_login").Update(user)
		if err != nil {
			log.Printf("更新登录时间失败: %v", err)
		}
		return user, true
	}
	
	return nil, false
}