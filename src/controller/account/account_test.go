package account

import (
	"crypto/md5"
	"fmt"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/utils"
	"net/http"
	"sync"
	"testing"
)

func TestAccount(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var req putPassword1Request
		req.Origin = fmt.Sprintf("%x", md5.Sum([]byte("admin")))
		req.New = fmt.Sprintf("%x", md5.Sum([]byte("123456")))
		code, _, _, err := utils.Request(http.MethodPut, "http://127.0.0.1:8088/api/1/password/admin", nil, &req, "")
		if err != nil || code != http.StatusOK {
			t.Error("Failed to PUT /api/1/password/admin. error: ", err)
			return
		}
		code, _, _, err = utils.Request(http.MethodPut, "http://127.0.0.1:8088/api/1/password/admin", nil, &req, "")
		if err != nil {
			t.Error("Failed to PUT /api/1/password/admin. error: ", err)
		}
		if code == http.StatusOK {
			t.Error("Changing password should not be success.")
		}
		temp := req.Origin
		req.Origin = req.New
		req.New = temp
		code, _, _, err = utils.Request(http.MethodPut, "http://127.0.0.1:8088/api/1/password/admin", nil, &req, "")
		if err != nil || code != http.StatusOK {
			t.Error("Failed to PUT /api/1/password/admin. error: ", err)
		}
	}()

	go func() {
		wg.Wait()
		_ = httpserver.Shutdown()
	}()

	_ = httpserver.ListenAndServe()
}
