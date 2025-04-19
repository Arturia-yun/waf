package models

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"

	"Waf-control/modules/util"
	"Waf-control/setting"
	"path/filepath"
)

var Engine *xorm.Engine

// Init 初始化数据库连接
func Init() {
	var err error
	
	switch setting.DBType {
	case "mysql":
		Engine, err = xorm.NewEngine(setting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			setting.DBUser, setting.DBPassword, setting.DBHost, setting.DBName))
	case "sqlite3":
		// 确保数据库目录存在
		dbDir := filepath.Dir(setting.DBPath)
		os.MkdirAll(dbDir, os.ModePerm)
		
		Engine, err = xorm.NewEngine(setting.DBType, setting.DBPath)
	default:
		log.Fatalf("不支持的数据库类型: %s", setting.DBType)
	}
	
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	
	// 设置数据库连接池
	Engine.SetMaxIdleConns(5)
	Engine.SetMaxOpenConns(10)
	
	// 同步数据结构
	err = Engine.Sync2(new(Rule), new(Site), new(User))
	if err != nil {
		log.Fatalf("数据库同步失败: %v", err)
	}
	
	// 检查是否需要创建默认用户
	count, err := Engine.Count(new(User))
	if err != nil {
		log.Fatalf("检查用户数量失败: %v", err)
	}
	
	if count == 0 {
		// 创建默认管理员用户
		admin := &User{
			Username: "admin",
			Password: util.MakeMd5("admin"),
		}
		_, err = Engine.Insert(admin)
		if err != nil {
			log.Fatalf("创建默认管理员用户失败: %v", err)
		}
		log.Println("已创建默认管理员用户: admin/admin")
	}
}
