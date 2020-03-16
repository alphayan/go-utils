package db

import utils "github.com/alphayan/go-utils"

type Role struct {
	utils.Model
	Name string `gorm:"unique_index"`
}
