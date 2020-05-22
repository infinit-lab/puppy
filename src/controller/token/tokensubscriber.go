package token

import (
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/log"
	"github.com/infinit-lab/taiji/src/model/token"
	"github.com/infinit-lab/yolanda/bus"
	"time"
)

type tokenSubscriber struct {
}

func (s *tokenSubscriber) Handle(key int, value *bus.Resource) {
	t, ok := value.Data.(*token.Token)
	if !ok {
		return
	}
	l := log.LoginLog{
		Username: t.Username,
		Ip:       t.Ip,
		Time:     time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	switch value.Status {
	case base.StatusCreated:
		l.IsLogin = true
	case base.StatusDeleted:
		l.IsLogin = false
	default:
		return
	}
	_ = log.CreateLoginLog(&l)
}
