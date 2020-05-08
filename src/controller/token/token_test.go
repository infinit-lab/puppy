package token

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/infinit-lab/taiji/src/model/token"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/utils"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestControllerToken(t *testing.T) {
	os.Args = append(os.Args, "token.duration=1")
	os.Args = append(os.Args, "token.checkDuration=3")
	os.Args = append(os.Args, "server.port=8088")
	config.Exec()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		testHandlePostToken1(t)
		testHandleDeleteToken1(t)
		testAutoDeleteToken(t)
	}()

	go func() {
		wg.Wait()
		_ = httpserver.Shutdown()
	}()

	_ = httpserver.ListenAndServe()
}

var myToken string

func testHandlePostToken1(t *testing.T) {
	var request postToken1Request
	request.Username = "admin"
	request.Password = fmt.Sprintf("%x", md5.Sum([]byte("admin")))
	code, rep, _, err := utils.Request(http.MethodPost, "http://127.0.0.1:8088/api/1/token", nil, request, "")
	if err != nil || code != http.StatusOK {
		t.Error("Failed to request POST /api/1/token. error: ", err)
		return
	}
	var response postToken1Response
	err = json.Unmarshal(rep, &response)
	if err != nil {
		t.Error("Failed to unmarshal response. error: ", err)
		return
	}
	myToken = response.Data
	temp, err := token.GetToken(myToken)
	if err != nil {
		t.Error("Failed to GetToken. error: ", err)
	}
	if temp.Username != "admin" {
		t.Errorf("Username is %s. it should be admin", temp.Username)
	}
}

func testHandleDeleteToken1(t *testing.T) {
	code, _, _, err := utils.Request(http.MethodDelete, "http://127.0.0.1:8088/api/1/token/"+myToken, nil, nil, "")
	if err != nil || code != http.StatusOK {
		t.Errorf("Failed to request DELETE /api/1/token/%s. error: %s", myToken, err)
		return
	}
	_, err = token.GetToken(myToken)
	if err == nil {
		t.Error("Should not get the token ", myToken)
	}
}

func testAutoDeleteToken(t *testing.T) {
	testHandlePostToken1(t)
	time.Sleep(5 * time.Second)
	_, err := token.GetToken(myToken)
	if err == nil {
		t.Error("Should not get the token ", myToken)
	}
}
