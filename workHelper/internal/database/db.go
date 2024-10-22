package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

//
//// 定义模型
//type Company struct {
//	ID      int     `gorm:"id"`
//	Name    string  `gorm:"name"`
//	Age     int     `gorm:"age"`
//	Address string  `gorm:"address"`
//	Salary  float64 `gorm:"salay"`
//}
//
//func main() {
//	db, err := gorm.Open(sqlite.Open("/Users/zhangrui/Workspace/goSpace/src/Tokumicn/theBookofChangesEveryDay/mhxyHelper/mhxyhelper.db"), &gorm.Config{})
//	if err != nil {
//		log.Fatal("failed to connect database")
//		return
//	}
//	db = db.Table("COMPANY")
//
//	companys := make([]Company, 0)
//	err = db.Model(Company{}).Debug().Find(&companys).Error
//	if err != nil {
//		log.Fatal("failed to connect database")
//		return
//	}
//
//	log.Println(companys)
//}

var (
	_db *gorm.DB
)

// 定义模型结构体
type User struct {
	gorm.Model
	Name string
	Age  uint8
}

// 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dir, _ := os.Getwd()
	fmt.Println("[MHXYDB] work dir: ", dir)

	db, err := gorm.Open(sqlite.Open("./mhxyhelper.db"), &gorm.Config{
		//// 开启 WAL 模式
		//DSN: "mode=wal",
		//// 增加最大连接数为 100
		//MaxOpenConns: 100,
	})
	if err != nil {
		return nil, err
	}
	// 设置数据库连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	// TODO 开启SQL打印测试用
	// db = db.Debug()

	_db = db

	return _db, nil
}

// 获取数据库连接
func LocalDB() *gorm.DB {
	if _db != nil {
		return _db
	}
	panic("database connection is nil")
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

// 初始化数据库并创建表
func InitDBWithAutoMigrate(needAutoMigrate bool) (*gorm.DB, error) {
	// 初始化数据库连接
	db, err := InitDB()
	if err != nil {
		panic(err)
	}

	if needAutoMigrate {

		// 物品信息表
		db.AutoMigrate(Stuff{})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		// 物品更新日志
		db.AutoMigrate(StuffLog{})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

	}

	return db, nil
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
