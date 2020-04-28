package token

import (
	"github.com/infinit-lab/puppy/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/logutils"
	"testing"
	"time"
)

var token string

type tokenHandler struct {

}

func (h* tokenHandler) Handle(key int, resource *bus.Resource) {
	switch key {
	case base.KeyToken:
		token := resource.Data.(*Token)
		logutils.TraceF("Token is %s, status is %d, username is %s, duration is %d, ip is %s, time is %s",
			resource.Id, resource.Status, token.Username, token.Duration, token.Ip, token.Time)
	default:
		logutils.Error("Should not receive key ", key)
	}
}

func TestNotification(t *testing.T) {
	bus.Subscribe(base.KeyToken, new(tokenHandler))
}

func TestCreateToken(t *testing.T) {
	var err error
	token, err = CreateToken("admin", 30, "127.0.0.1", nil)
	if err != nil {
		t.Error("Failed to CreateToken. error: ", err)
	}
	temp, err := GetToken(token)
	if err != nil {
		t.Error("Failed to GetToken. error: ", err)
	}
	if temp.Username != "admin" {
		t.Errorf("Username is %s. it should be admin", temp.Username)
	}
	if temp.Duration != 30 {
		t.Errorf("Duration is %d, it should be 30", temp.Duration)
	}
	if temp.Ip != "127.0.0.1" {
		t.Errorf("Ip is %s, it should be 127.0.0.1", temp.Ip)
	}
}

func TestRenewToken(t *testing.T) {
	temp, err := GetToken(token)
	if err != nil {
		t.Error("Failed to GetToken. error: ", err)
		return
	}
	time.Sleep(time.Second)
	err = RenewToken(token)
	if err != nil {
		t.Error("Failed to RenewToken. error: ", err)
	}
	temp2, err := GetToken(token)
	if err != nil {
		t.Error("Failed to GetToken. error: ", err)
	}
	if temp.Time == temp2.Time {
		t.Error("The time of token should not be same.")
	}
}

func TestGetTokenList(t *testing.T) {
	tokenList, err := GetTokenList()
	if err != nil {
		t.Error("Failed to GetTokenList. error: ", err)
	}
	if len(tokenList) == 0 {
		t.Error("Token list should not be empty")
	}
}

func TestDeleteToken(t *testing.T) {
	err := DeleteToken(token, nil)
	if err != nil {
		t.Error("Failed to DeleteToken. error: ", err)
	}
	_, err = GetToken(token)
	if err == nil {
		t.Error("The token should be deleted.")
	} else {
		logutils.Trace("Delete result is ", err)
	}
	time.Sleep(100 * time.Millisecond)
}
