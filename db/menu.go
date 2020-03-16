package db

import "github.com/jinzhu/gorm"

//import "github.com/alphayan/go-utils"

//type Menu struct {
//	utils.Model
//	UserID uint32
//	Name   string
//	Url    string
//	Order  int32 //排序
//	ParentID uint32
//	RootID uint32 //根目录id
//	Other  string
//}
//
//func (Menu) TableName() string {
//	return "core_menu"
//}
//
//func NewMenu() (m *Menu) {
//	m = &Menu{Model: utils.Model{DB: _db}}
//	m.SetParent(m)
//	return
//}

func MenuWhereArgs(id uint32) []interface{} {
	args := make([]interface{}, 3)
	args[0] = "type = ? and id = ?"
	args[1] = KV_TYPE_MENU
	args[2] = id
	return args
}

func MenuQuery(id uint32) *gorm.DB {
	menu := NewKV()
	return menu.Where(MenuWhereArgs(id))
}
func GetMenuById(id uint32) (*KV, error) {
	menu := NewKV()
	err := MenuQuery(id).Find(menu).Error
	return menu, err
}

func GetMenuByKey(key string) (*KV, error) {
	menu := NewKV()
	menu.Type = KV_TYPE_MENU
	menu.Key = key
	err := menu.Where(menu).Find(menu).Error
	return menu, err
}

func CheckMenuExist(id uint32) bool {
	count := 0
	MenuQuery(id).Count(&count)
	return count == 1
}
