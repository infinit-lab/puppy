package net

import (
	"encoding/json"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"testing"
)

func TestUpdateAddress(t *testing.T) {
	os.Args = append(os.Args, "log.level=3")
	config.Exec()
	var addr Address
	addr.Name = "1234"
	addr.Ip = "192.168.1.100"
	addr.Mask = "255.255.255.0"
	addr.Gateway = "192.168.1.1"
	err := UpdateAddress(&addr)
	if err != nil {
		t.Fatal("Failed to Update Address")
	}
	addrs, err := GetAddressList()
	if err != nil {
		t.Fatal("Failed to GetAddressList. error: ", err)
	}

	for _, a := range addrs {
		data, _ := json.Marshal(a)
		logutils.Trace("Address is ", string(data))
		a.Ip = "192.168.2.100"
		a.Gateway = "192.168.2.1"
		err := UpdateAddress(a)
		if err != nil {
			t.Fatal("Failed to UpdateAddress. error: ", err)
		}
		temp, err := GetAddress(a.Name)
		if err != nil {
			t.Fatal("Failed to GetAddress. error: ", err)
		}
		data, _ = json.Marshal(temp)
		logutils.Trace("Address is ", string(data))

		err = DeleteAddress(temp.Name)
		if err != nil {
			t.Fatal("Failed to DeleteAddress. error: ", err)
		}
	}
}
