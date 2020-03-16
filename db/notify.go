package db

import (
	"bytes"
	"fmt"
	"time"

	utils "github.com/alphayan/go-utils"
)

const (
	NotifyTypeNormal    uint32 = 0
	REQUIREMENT_TYPE    uint32 = 1 // 小车匠、主机厂推送需求的消息（教务端）
	STUDENT_RESUME_TYPE uint32 = 2 // 学生向教师提交个人简历的消息（教师端）
	CLASS_RESUME_TYPE   uint32 = 3 // 教师向教务推送班级学生简历的消息（教务端）
	CLASS_TYPE          uint32 = 4 // 主机厂处理了教务端提交的开班申请的消息（教务端、教师端）
	TEACHER_RESUME_TYPE uint32 = 5 // 教师向教务提交个人简历的消息（教务端）
	ExamType            uint32 = 6 //考试消息
	CertificateType     uint32 = 7 //证书到期消息
)

type NotifyMsg struct {
	utils.Model
	utils.ModelTime
	ToUserID uint32
	Context  string `gorm:"type:text"`
	Type     uint32 `gorm:"type:tinyint"`
	IsRead   bool
}

func (NotifyMsg) TableName() string {
	return "core_notify_msg"
}

func NewNotifyMsg() (m *NotifyMsg) {
	m = &NotifyMsg{Model: utils.Model{DB: _db}}
	m.SetParent(m)
	return
}

func (m *NotifyMsg) SendTo(uids []uint32) (err error) {
	if len(uids) == 0 {
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO core_notify_msg (created_at, context, type, is_read, to_user_id) VALUES ")
	sqlVal := fmt.Sprintf("('%s' , '%s', '%d', 0, ", time.Now().Format("2006-01-02 15:04:05"), utils.SqlEscape(m.Context), m.Type)
	for i, uid := range uids {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(sqlVal)
		buffer.WriteString(fmt.Sprintf("'%d')", uid))
	}
	err = m.Table().Exec(buffer.String()).Error
	return
}

// 保存消息
//ids：消息的接收者；msgType：消息类型；userId：教师的id
func NotifySaveInfo(ids []uint32, userId uint32, Title string, msgType uint32) error {
	var teacher Account
	if userId != 0 { // 查询教师的信息
		GetDB().Raw("SELECT acc.* FROM accounts acc WHERE acc.group = ? and acc.id = ?",
			GROUP_TEACHER, userId).
			Scan(&teacher)
	}
	for _, toUserID := range ids {
		msg := NewNotifyMsg()
		msg.ToUserID = toUserID
		msg.Type = msgType
		if msgType == REQUIREMENT_TYPE { // 小车匠、主机厂推送需求的消息（教务端）type=1
			msg.Context = "主机厂" + Title + ",推送了需求"
		} else if msgType == STUDENT_RESUME_TYPE { // 学生向教师提交个人简历的消息（教师端）type=2
			msg.Context = "学生：" + Title + "，提交了个人简历"
		} else if msgType == CLASS_RESUME_TYPE { // 教师向教务推送班级学生简历的消息（教务端）type=3
			msg.Context = "教师：" + teacher.Nickname + "，推送了班级学生的简历到班级：" + Title
		} else if msgType == CLASS_TYPE { // 主机厂处理了教务端提交的开班申请的消息（教务端）type=4
			msg.Context = "班级：" + Title + "，提交的开班申请已处理"
		} else {
			msg.Context = Title
		}
		err := msg.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
