package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	_db        *gorm.DB
	CanMigrate bool
)

func init() {
	//conf := vault.GetSecretValues("mysql")
	//db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	//	conf["user"], conf["passwd"], conf["host"], conf["port"], conf["db"]))
	//if err != nil {
	//	panic(err)
	//}
	//_db = db
	//db.AutoMigrate(&Account{})
}

func Connect(addr string, debug bool) *gorm.DB {
	db, err := gorm.Open("mysql", addr)
	if err != nil {
		fmt.Print("mysql addr:", addr)
		panic(err)
	}
	db.LogMode(debug)
	return db
}

func ConnectMain(addr string, debug bool, migrate bool) {
	_db = Connect(addr, debug)
	CanMigrate = migrate
	if CanMigrate {
		_db.AutoMigrate(&Account{})
		_db.AutoMigrate(&Role{})
		_db.AutoMigrate(&Media{})
		_db.AutoMigrate(&KV{})
		_db.AutoMigrate(&NotifyMsg{})
		//_db.AutoMigrate(&Menu{})

		NewNotifyMsg().Table().AddForeignKey("to_user_id", "accounts(id)", "CASCADE", "NO ACTION")
	}
}

func GetDB() *gorm.DB {
	return _db
}
