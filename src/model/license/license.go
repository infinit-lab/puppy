package license

import (
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"sync"
)

type License struct {
	Auth map[string]Auth
}

type Auth struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	ValueType string   `json:"valueType"` //int, bool, datetime, string
	Value     []string `json:"value"`
	Current   string   `json:"current,omitempty"`
}

var status int
var mutex sync.Mutex

func init() {
	status = base.LicenseUnauthorized
}

func SetLicenseStatus(s int) {
	if status != s {
		mutex.Lock()
		status = s
		mutex.Unlock()
		_ = bus.PublishResource(base.KeyLicenseStatus, base.StatusUpdated, "", status, nil)
	}
}

func GetLicenseStatus() int {
	mutex.Lock()
	defer mutex.Unlock()
	return status
}
