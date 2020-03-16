package db

import (
	"errors"

	utils "github.com/alphayan/go-utils"
)

const (
	KV_TYPE_MENU    = "MENU"     //菜单类型
	KV_TYPE_ABOUTUS = "ABOUT_US" //关于我们
)

type KV struct {
	utils.Model
	Key       string `gorm:"type:varchar(200);not null;unique_index"`
	Name      string `gorm:"type:varchar(200)"` //描述
	UserID    uint32
	Value     string `gorm:"type:text"`
	Type      string `gorm:"index"` //特殊类型标志
	UserGroup uint32 //有些kv需要访问权限
}

func (KV) TableName() string {
	return "core_kv"
}

func NewKV() (m *KV) {
	m = &KV{Model: utils.Model{DB: _db}}
	m.SetParent(m)
	return
}

func (m *KV) UpdateValue(v string) error {
	return m.Update("value", v)
}

func (m *KV) UpdateNameValue(name, value string) error {
	return m.Updates(map[string]interface{}{"value": value, "name": name})
}

func (p *KV) IsValid() error {
	errs := make([]error, 1)
	if p.Key == "" {
		errs = append(errs, errors.New("key is empty"))
	}
	return utils.FirstError(errs...)
}

func GetKvByKey(key string) (*KV, error) {
	menu := NewKV()
	menu.Key = key
	err := menu.Where(menu).Find(menu).Error
	return menu, err
}
