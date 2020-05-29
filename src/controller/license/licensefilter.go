package license

import (
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/license"
	"github.com/infinit-lab/yolanda/httpserver"
	"net/http"
)

type licenseFilter struct {
}

func (f *licenseFilter) Filter(r *http.Request, checkToken bool) error {
	if checkToken == false {
		return nil
	}
	if license.GetLicenseStatus() == base.LicenseAuthorized {
		return nil
	}
	if r.Method == http.MethodGet {
		return nil
	}
	if license.GetLicenseStatus() == base.LicenseImporting {
		return errors.New("正在导入授权，请稍后重试")
	}
	return errors.New("未授权")
}

var filter licenseFilter

func init() {
	httpserver.RegisterFilter(&filter)
}
