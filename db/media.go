package db

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/alphayan/go-utils"
	"github.com/rs/xid"
	"gopkg.in/h2non/filetype.v1"
)

const (
	UPLOAD_URL_PREFIX              = "/uploads"
	MEDIA_TYPE_UNKNOW              = 0
	MEDIA_TYPE_IMAGE               = 1
	MEDIA_TYPE_VIDEO               = 2
	MEDIA_TYPE_EXCEL               = 15
	MEDIA_TYPE_PX                  = 10 //自定义类型
	MEDIA_TYPE_YUNCOURSE           = 11 //自定义类型,云课堂
	MEDIA_TYPE_WEXAM_XcjSimulation = 20 //考试，仿真题
)

//设置默认路径
var MediaPath string = "temp"

type Media struct {
	utils.Model
	utils.ModelTime
	UserID          uint32 //上传者
	LocalPath       string `gorm:"index"`
	Name            string
	Description     string `grom:"type:varchar(100)"`
	FileExt         string //.png .jpg
	FileType        int32
	ParentID        uint32 `gorm:"index"`
	XcjQuestionId   uint32 //小车匠试题id
	XcjCourseId     uint32
	XcjCourseNumber string `gorm:"type:varchar(255)"` //记录下当前课件的课程编号
}

func (Media) TableName() string {
	return "cms_media"
}

func NewMedia() (m *Media) {
	m = &Media{Model: utils.Model{DB: _db}}
	m.SetParent(m)
	return
}

func (p *Media) SaveFile(file multipart.File) error {
	guid := xid.New()
	fileSubPath := guid.Time().Format("2006-01")
	fileDir := filepath.Join(MediaPath, fileSubPath)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, os.ModePerm)
	}
	fileLocalName := fmt.Sprintf("%s%s", guid.String(), p.FileExt)
	out, err := os.OpenFile(filepath.Join(fileDir, fileLocalName),
		os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	p.LocalPath = path.Join(fileSubPath, fileLocalName)
	buf := make([]byte, 100)
	file.Seek(0, 0)
	file.Read(buf)
	if filetype.IsImage(buf) {
		p.FileType = MEDIA_TYPE_IMAGE
	} else {
		p.FileType = MEDIA_TYPE_UNKNOW
	}
	err = p.Save()
	return err
}

func (p Media) FullPath() string {
	return GetMediaFullUrl(p.LocalPath)
}

func (p Media) GetFullLocalPath() string {
	return GetMediaFullPath(p.LocalPath)
}

func GetMediaFullUrl(localpath string) string {
	if !utils.StrIsEmpty(localpath) {
		return "/uploads/" + localpath
	}
	return ""
}
func GetMediaFullPath(localpath string) string {
	return path.Join(MediaPath, localpath)
}

func GetMediaFullPathWithMonth(filename string) string {
	fileSubPath := time.Now().Format("2006-01")
	fileDir := filepath.Join(MediaPath, fileSubPath)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, 0771)
	}
	return filepath.Join(fileDir, filename)
}
