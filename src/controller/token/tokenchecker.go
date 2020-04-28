package token

import (
	"errors"
	"github.com/infinit-lab/puppy/src/model/token"
	"net/http"
)

type tokenChecker struct {
}

func (c *tokenChecker) CheckToken(r *http.Request) error {
	auth, ok := r.Header["Authorization"]
	if !ok || len(auth) == 0 {
		return errors.New("无效Token")
	}
	_, err := token.GetToken(auth[0])
	if err != nil {
		return err
	}
	_ = token.RenewToken(auth[0])
	return nil
}
