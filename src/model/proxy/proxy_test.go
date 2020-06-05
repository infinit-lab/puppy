package proxy

import (
	"encoding/json"
	"github.com/infinit-lab/qiankun/common"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	uuid "github.com/satori/go.uuid"
	"os"
	"strings"
	"testing"
)

var s common.Server

func checkServer(t *testing.T) {
	server, err := GetLocalServer(s.Uuid)
	if err != nil {
		t.Fatal("Failed to GetLocalServer. error: ", err)
	}
	if server.Uuid != s.Uuid || server.Host != s.Host ||
		server.Port != s.Port || server.Description != s.Description {
		t.Fatal("The data of local server is wrong.")
	}
}

func TestCreateLocalServer(t *testing.T) {
	s.Uuid = strings.ReplaceAll(uuid.NewV1().String(), "-", "")
	s.Host = "127.0.0.1"
	s.Port = 8080
	s.Description = "test local server"
	err := CreateLocalServer(&s)
	if err != nil {
		t.Fatal("Failed to CreateLocalServer. error: ", err)
	}
	checkServer(t)
}

func TestUpdateLocalServer(t *testing.T) {
	s.Host = "172.21.0.1"
	s.Port = 8088
	s.Description = "update local server"
	err := UpdateLocalServer(&s)
	if err != nil {
		t.Fatal("Failed to UpdateLocalServer. error: ", err)
	}
	checkServer(t)
}

func TestGetLocalServerList(t *testing.T) {
	os.Args = append(os.Args, "log.level=3")
	config.Exec()
	servers, err := GetLocalServerList()
	if err != nil {
		t.Fatal("Failed to GetLocalServerList. error: ", err)
	}
	for _, server := range servers {
		data, _ := json.Marshal(server)
		logutils.Trace("Server is ", string(data))
	}
}

func TestDeleteLocalServer(t *testing.T) {
	err := DeleteLocalServer(s.Uuid)
	if err != nil {
		t.Fatal("Failed to DeleteLocalServer. error: ", err)
	}
	_, err = GetLocalServer(s.Uuid)
	if err == nil {
		t.Fatal("Should not get local server.")
	}
}

func TestCreateRemoteHost(t *testing.T) {
	r := RemoteHost {
		Address: "127.0.0.1:7070",
		Description: "Test Remote Host",
	}
	err := CreateRemoteHost(&r)
	if err != nil {
		t.Fatal("Failed to CreateRemoteHost. error: ", err)
	}
	rs, err := GetRemoteHostList()
	if err != nil {
		t.Fatal("Failed to GetRemoteHostList. error: ", err)
	}
	for _, r := range rs {
		data, _ := json.Marshal(r)
		logutils.Trace("Remote host is ", string(data))
	}
}

func TestDeleteRemoteHost(t *testing.T) {
	err := DeleteRemoteHost("127.0.0.1:7070")
	if err != nil {
		t.Fatal("Failed to DeleteRemoteHost. error: ", err)
	}
	rs, err := GetRemoteHostList()
	for _, r := range rs {
		data, _ := json.Marshal(r)
		logutils.Trace("Remote host is ", string(data))
	}
}
