package db

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	utils "github.com/alphayan/go-utils"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	"github.com/alphayan/iris"
	funk "github.com/thoas/go-funk"
	"golang.org/x/crypto/pbkdf2"
)

type AccountType int

const (
	TYPE_PLAIN   AccountType = 1
	TYPE_WX      AccountType = 2
	GROUP_BANNED uint32      = 0 //禁用
	GROUP_ADMIN  uint32      = 1
	GROUP_USER   uint32      = 20
	//该项目新增组
	GROUP_TEACHER uint32 = 21
	GROUP_STUDENT uint32 = 22
)

var (
	AccountGroupList           = []uint32{GROUP_ADMIN, GROUP_USER, GROUP_TEACHER, GROUP_STUDENT}
	DefaultPassword     string = "123456"
	ErrPasswordNotMatch        = errors.New("密码错误")
)

type Account struct {
	utils.Model
	utils.ModelTime
	Username string `gorm:"type:varchar(30);unique_index"`
	//Mail     string `gorm:"type:varchar(100);unique_index"`
	//Phone    string `gorm:"type:varchar(20);unique_index"`
	//当前项目不需mail和phone
	Mail     string `gorm:"type:varchar(100);"`
	Phone    string `gorm:"type:varchar(20);"`
	Salt     string `json:"-"`
	Password string `json:"-"`
	Nickname string
	Type     AccountType
	//用户组
	Group uint32 `gorm:"type:tinyint"`
	//Role   Role
	//RoleID uint32
}

func NewAccount() (m *Account) {
	m = &Account{Model: utils.Model{DB: _db}}
	fmt.Println(_db)
	m.SetParent(m)
	return
}

// 创建用户
// 使用 32bit的salt
// key := pbkdf2.Key([]byte(password), salt, 1024, 32, sha1.New)
// skey := base64.StdEncoding.EncodeToString(key)
// 最后保存base64的密码字符串
func (acc *Account) Create() error {
	if acc.Salt != "" {
		return errors.New("account already has a salt")
	}
	//默认为用户组
	if acc.Group == 0 {
		acc.Group = GROUP_USER
	}
	switch acc.Type {
	case TYPE_PLAIN:
		skey, salt := CreatePassword(acc.Password)
		acc.Salt = salt
		acc.Password = skey
	default:
		return errors.New(fmt.Sprintf("account no such type %v", acc.Type))
	}
	err := acc.Save()
	return err
}

func (acc *Account) BeforeCreate() (err error) {
	if !funk.Contains(AccountGroupList, acc.Group) {
		return errors.New("need group")
	}
	if acc.Username == "" {
		err = errors.New("username is empty")
	}
	return err
}

func (acc *Account) FormatError(err error) error {
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "uix_accounts_username") {
			err = errors.New("username existed")
		} else if strings.Contains(errStr, "uix_accounts_mail") {
			err = errors.New("mail重复")
		} else if errStr == "record not found" {
			err = errors.New("账号不存在")
		} else {
			err = acc.Model.FormatError(err)
		}
	}
	return err
}
func (acc *Account) CheckPassword(password string) bool {
	acc.FetchColumnValue("Password", "Salt")
	salt, _ := base64.StdEncoding.DecodeString(acc.Salt)
	key := pbkdf2.Key([]byte(password), salt, 1024, 32, sha256.New)

	if base64.StdEncoding.EncodeToString(key) == acc.Password {
		return true
	}
	return false
}

// 创建 pbkdf2的密码
// skey, slats base64后的字符串
func CreatePassword(password string) (skey string, salts string) {
	salt := make([]byte, 32)
	_, err := crand.Reader.Read(salt)
	if err != nil {
		//todo error
	}
	//使用 pbkdf2 生产秘钥
	key := pbkdf2.Key([]byte(password), salt, 1024, 32, sha256.New)
	skey = base64.StdEncoding.EncodeToString(key)
	salts = base64.StdEncoding.EncodeToString(salt)
	return
}

// 验证用户名和密码
// 返回nil时，用户名或密码错误
func CheckAccountPassword(name string, password string) (*Account, error) {
	acc, err := GetAccountByName(name, TYPE_PLAIN)
	if err != nil {
		return nil, err
	}
	if acc.CheckPassword(password) {
		return acc, nil
	}
	return nil, ErrPasswordNotMatch
}

func GetAccountByName(name string, atype AccountType) (acc *Account, err error) {
	acc = NewAccount()
	f := _db.Where("type = ?", atype)
	switch atype {
	case TYPE_PLAIN:
		if strings.Contains(name, "@") {
			err = f.Where("mail = ?", name).First(&acc).Error
		} else {
			//TODO 查询手机号
			err = f.Where("username = ?", name).First(&acc).Error
		}
	default:
		err = errors.New("no such account type")
	}
	err = acc.FormatError(err)
	return
}

//判断账号是否可用
func IsAccountAvailable(uid uint32) bool {
	var c int
	_db.Model(&Account{}).Where("id = ?", uid).Count(&c)
	return c > 0
}

func CheckGroup(needGroup, thisGroup uint32) error {
	//通用权限
	if needGroup == 0 || needGroup == GROUP_USER ||
		thisGroup == GROUP_ADMIN || //管理员有完全访问权限
		needGroup == thisGroup {
		return nil
	}

	// TODO 权限比较
	return errors.New("no permission")
}

func JoinAccountName(g *gorm.DB, mod utils.IModel, selectKey, joinKey string) *gorm.DB {
	attrs := utils.GetSelectAttrs(g)
	attrs = append(attrs, fmt.Sprintf("accounts.nickname as %v", selectKey))
	g = g.Select(attrs).Joins(
		fmt.Sprintf("LEFT JOIN accounts ON accounts.id = %s.%s", mod.GetTableName(), joinKey))
	return g
}

// 管理员列表
func GetAdminList() map[string]interface{} {
	admins := make([]*Account, 0)
	err := GetDB().Raw("SELECT acc.* FROM accounts acc WHERE acc.group =  ?", GROUP_ADMIN).
		Scan(&admins).Error
	if err != nil {
		return map[string]interface{}{"State": iris.StatusBadRequest, "Result": "查询管理员出错", "Err": err}
	} else {
		var ids []uint32
		for _, admin := range admins {
			ids = append(ids, admin.ID)
		}
		return map[string]interface{}{"State": iris.StatusOK, "Result": ids}
	}
}
