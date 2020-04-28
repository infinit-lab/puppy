package token

import (
	"errors"
	"github.com/infinit-lab/puppy/src/model/token"
	"net/http"
	"net/url"
)

type tokenChecker struct {
}

func (c *tokenChecker) CheckToken(r *http.Request) error {
	auth, ok := r.Header["Authorization"]
	if !ok || len(auth) == 0 {
		form, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return errors.New("无效Token")
		}
		auth, ok = form["token"]
		if !ok || len(auth) == 0 {
			return errors.New("无效Token")
		}
	}
	_, err := token.GetToken(auth[0])
	if err != nil {
		return err
	}
	_ = token.RenewToken(auth[0])
	return nil
}
