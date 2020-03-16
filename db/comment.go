package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

//评分系统
type CommentR struct {
	//utils.Model
	CommentRefID uint32     //对应的需要评分的资源 comment_ref_id
	FromUserID   uint32     //评论人 from_user_id
	CommentTime  *time.Time //评论时间
	CommentText  string     `gorm:"type:varchar(200)"` //评论留言
	CommentScore uint32     //评分
}

type CommentAvgR struct {
	CommentRefID    uint32
	CommentCount    uint32  //comment_count
	CommentAvgScore float32 //comment_avg_score
}

func (m *CommentR) GetCommentScoreAvg(g *gorm.DB) float32 {
	av := CommentAvgR{}
	g.Select("AVG(comment_score) as comment_avg_score").Scan(&av)
	return av.CommentAvgScore
}

func (m *CommentR) IsUserCommented(g *gorm.DB) bool {
	ac := CommentR{}
	err := g.Select("from_user_id").Where("from_user_id = ? AND comment_ref_id = ?",
		m.FromUserID, m.CommentRefID).Limit(1).Scan(&ac).Error
	return err == nil && ac.FromUserID > 0
}
