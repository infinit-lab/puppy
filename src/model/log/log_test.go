package log

import (
	"encoding/json"
	"github.com/infinit-lab/yolanda/logutils"
	"testing"
	"time"
)

func TestCreateOperateLog(t *testing.T) {
	l := OperateLog{
		Username:    "admin",
		Ip:          "127.0.0.1",
		Operate:     "started",
		ProcessId:   3,
		ProcessName: "ping",
		Time:        time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	err := CreateOperateLog(&l)
	if err != nil {
		t.Fatal("Failed to CreateOperateLog. error: ", err)
	}
}

func TestGetOperateLogList(t *testing.T) {
	date := time.Now().Local().Format("2006-01-02")
	logList, err := GetOperateLogList(date+" 00:00:00", date+" 23:59:59", "", 0, 100, 0)
	if err != nil {
		t.Fatal("Failed to GetOperateLogList. error: ", err)
	}
	if len(logList) == 0 {
		t.Fatal("Log count should not be empty")
	}
	for _, l := range logList {
		data, err := json.Marshal(l)
		if err != nil {
			t.Fatal("Failed to Marshal. error: ", err)
		}
		logutils.Error("log is ", string(data))
	}
}

func TestCreateLoginLog(t *testing.T) {
	l := LoginLog{
		Username: "admin",
		Ip:       "127.0.0.1",
		IsLogin:  true,
		Time:     time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	err := CreateLoginLog(&l)
	if err != nil {
		t.Fatal("Failed to CreateLoginLog. error: ", err)
	}
}

func TestGetLoginLogList(t *testing.T) {
	date := time.Now().Local().Format("2006-01-02")
	logList, err := GetLoginLogList(date+" 00:00:00", date+" 23:59:59", 100, 0)
	if err != nil {
		t.Fatal("Failed to GetLoginLogList. error: ", err)
	}
	if len(logList) == 0 {
		t.Fatal("Log count should not be empty")
	}
	for _, l := range logList {
		data, err := json.Marshal(l)
		if err != nil {
			t.Fatal("Failed to Marshal. error: ", err)
		}
		logutils.Error("log is ", string(data))
	}
}
