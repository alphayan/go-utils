package uhandler

import (
	. "github.com/alphayan/go-utils"
	cdb "github.com/alphayan/go-utils/db"
	. "github.com/alphayan/go-utils/session"
	"github.com/alphayan/iris/context"
)

func RequireLogin(ictx context.Context) {
	sess, _ := NewSessionFromIris(ictx, XSESSION_KEY)
	RequireLoginX(sess)
}

func RequireLoginX(sess *XSession) {
	if sess.Uid <= 0 {
		panic(ERR_REQUIRE_LOGIN)
	}
}

func RequireAdmin(ictx context.Context) {
	sess, _ := NewSessionFromIris(ictx, XSESSION_KEY)
	if !sess.IsAdmin() {
		panic(ERR_REQUIRE_ADMIN)
	}
}

//PX
func RequireTeacher(sess *XSession) {
	err := cdb.CheckGroup(cdb.GROUP_TEACHER, sess.Group)
	if err != nil {
		panic("need teacher")
	}
}
