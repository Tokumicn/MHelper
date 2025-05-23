package data

import (
	"fmt"
	"gorm.io/gorm"
)

// 定义模型结构体
type User struct {
	gorm.Model
	Name string
	Age  uint8
}

// 定义批量写入函数
func batchInsertUsers(db *gorm.DB, users []User) error {
	// 每次写入 1000 条数据
	batchSize := 1000
	batchCount := (len(users) + batchSize - 1) / batchSize
	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(users) {
			end = len(users)
		}
		batch := users[start:end]
		// 启用事务
		tx := db.Begin()
		if err := tx.Error; err != nil {
			return err
		}
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

// 查询用户信息
func getUsers(db *gorm.DB) ([]User, error) {
	var users []User
	// 使用缓存，减少对数据库的读操作
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func createUserDatas(db *gorm.DB) {
	// 批量插入数据
	users := []User{}
	for i := 0; i < 1000; i++ {
		user := User{
			Name: "user_" + string(i),
			Age:  uint8(i % 100),
		}
		users = append(users, user)
	}
	err := batchInsertUsers(db, users)
	if err != nil {
		panic(err)
	}

	// 查询数据
	users, err = getUsers(db)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}
