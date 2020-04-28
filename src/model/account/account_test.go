package account

import (
	"testing"
)

func TestInit(t *testing.T) {

}

func TestIsValidAccount(t *testing.T) {
	isValid, err := IsValidAccount("admin", "admin")
	if err != nil {
		t.Error("Failed to IsValidAccount. error: ", err)
		return
	}
	if isValid == true {
		t.Error("Should not be valid account")
	}

	isValid, err = IsValidAccount("admin", "21232f297a57a5a743894a0e4a801fc3")
	if err != nil {
		t.Error("Failed to IsValidAccount. error: ", err)
		return
	}
	if isValid == false {
		t.Error("Should be valid account")
	}
}

func TestChangePassword(t *testing.T) {
	err := ChangePassword("admin", "21232f297a57a5a7", "admin", nil)
	if err == nil {
		t.Error("ChangePassword should no be success")
	}
	err = ChangePassword("admin", "21232f297a57a5a743894a0e4a801fc3", "admin", nil)
	if err != nil {
		t.Error("Failed to ChangePassword. error: ", err)
	}
	isValid, err := IsValidAccount("admin", "admin")
	if err != nil {
		t.Error("Failed to IsValidAccount. error: ", err)
	}
	if isValid == false {
		t.Error("Should be valid account")
	}
	err = ChangePassword("admin", "admin", "21232f297a57a5a743894a0e4a801fc3", nil)
	if err != nil {
		t.Error("Failed to ChangePassword. error: ", err)
	}
	isValid, err = IsValidAccount("admin", "21232f297a57a5a743894a0e4a801fc3")
	if err != nil {
		t.Error("Failed to IsValidAccount. error: ", err)
	}
	if isValid == false {
		t.Error("Should be valid account")
	}
}
