package token

import (
	"github.com/infinit-lab/puppy/src/model/base"
	"github.com/infinit-lab/puppy/src/model/token"
	"github.com/infinit-lab/yolanda/bus"
)

type passwordSubscriber struct {
}

func (h *passwordSubscriber) Handle(key int, value *bus.Resource) {
	if key != base.KeyPassword || value.Status != base.StatusUpdated {
		return
	}

	username, ok := value.Data.(string)
	if !ok {
		return
	}

	tokenList, err := token.GetTokenList()
	if err != nil {
		return
	}

	for _, t := range tokenList {
		if t.Username == username {
			_ = token.DeleteToken(t.Token, value.Context)
		}
	}
}
