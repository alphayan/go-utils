package db

import (
	"time"

	utils "github.com/alphayan/go-utils"
	"github.com/jinzhu/gorm"
)

//审核功能
const (
	REVIEW_STATE_WAIT    = 0 //待审核
	REVIEW_STATE_SUCCESS = 1 //审核通过
	REVIEW_STATE_FAIL    = 2 //审核未通过
)

var (
	ReviewStateList   = []uint32{REVIEW_STATE_WAIT, REVIEW_STATE_SUCCESS, REVIEW_STATE_FAIL}
	ReviewStateString = map[uint32]string{
		REVIEW_STATE_WAIT:    "未审核",
		REVIEW_STATE_SUCCESS: "审核通过",
		REVIEW_STATE_FAIL:    "审核未通过",
	}
)

type Review struct {
	//utils.Model
	RefID        uint32     //对应的需要审核的资源
	FromUserID   uint32     //资源来源方 from_user_id
	FromUserName string     `gorm:"-"`
	ReviewState  uint32     `gorm:"type:tinyint"` //审核状态
	ReviewUserID uint32     //审核人
	ReviewTime   *time.Time //盛和时间
	ReviewText   string     `gorm:"type:varchar(200)"` //审核留言
}

func (m *Review) UpdateReviewState(db *gorm.DB, item map[string]interface{}) (err error) {
	//自动更新到现在
	item["ReviewTime"] = time.Now()
	err = db.Select("ReviewState", "ReviewUserID", "ReviewText", "ReviewTime", "CategoryID").UpdateColumns(item).Error
	return
}

func GetReviewStateStr(t uint32) string {
	if r, ok := ReviewStateString[t]; ok {
		return r
	}
	return ""
}

func (m *Review) ResetReviewState() {
	m.ReviewState = REVIEW_STATE_WAIT
	m.ReviewUserID = 0
	m.ReviewTime = nil
	m.ReviewText = ""
}

func ValidReviewState(res uint32) bool {
	return utils.ContainsUint32(ReviewStateList, res)
}
