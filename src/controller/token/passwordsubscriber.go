package token

import (
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/token"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/logutils"
)

type passwordSubscriber struct {
}

func (h *passwordSubscriber) Handle(key int, value *bus.Resource) {
	if key != base.KeyPassword || value.Status != base.StatusUpdated {
		logutils.ErrorF("Key is %d, status is %d", key, value.Status)
		return
	}

	tokenList, err := token.GetTokenList()
	if err != nil {
		logutils.Error("Failed to GetTokenList. error: ", err)
		return
	}

	for _, t := range tokenList {
		if t.Username == value.Id {
			_ = token.DeleteToken(t.Token, value.Context)
		}
	}
}
