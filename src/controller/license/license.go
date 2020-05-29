package license

import (
	"encoding/json"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/license"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

var l license.License

func init() {
	logutils.Trace("Initializing controller license...")
	loadLicense()
	checkLicense()
	go updateLicenseLoop()
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/fingerprint", HandleGetFingerprint1, true)
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/license", HandleGetLicense1, true)
	httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/license", HandlePutLicense1, false)
}

type getFingerprint1Response struct {
	httpserver.ResponseBody
	Data string `json:"data"`
}

func HandleGetFingerprint1(w http.ResponseWriter, r *http.Request) {
	var response getFingerprint1Response
	response.Result = true
	var err error
	response.Data, err = utils.GetMachineFingerprint()
	if err != nil {
		logutils.Error("Failed to GetMachineFingerprint. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpserver.Response(w, response)
}

type getLicense1Response struct {
	httpserver.ResponseBody
	Data struct {
		Status        int    `json:"status"`
		Fingerprint   string `json:"fingerprint"`
		IsForever     bool   `json:"isForever"`
		ValidDatetime string `json:"valueDatetime"`
		ValidDuration int    `json:"validDuration"`
	} `json:"data"`
}

func HandleGetLicense1(w http.ResponseWriter, r *http.Request) {
	var response getLicense1Response
	var err error
	response.Data.Fingerprint, err = utils.GetMachineFingerprint()
	if err != nil {
		logutils.Error("Failed to GetMachineFingerprint. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Data.Status = license.GetLicenseStatus()
	auth := l.Auth
	for key, value := range auth {
		if len(value.Value) == 0 {
			continue
		}
		switch key {
		case base.AuthForever:
			response.Data.IsForever, err = strconv.ParseBool(value.Value[0])
			if err != nil {
				logutils.Error("Failed to ParseBool. error: ", err)
				httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case base.AuthDatetime:
			response.Data.ValidDatetime = value.Value[0]
		case base.AuthDuration:
			d, err := strconv.Atoi(value.Value[0])
			if err != nil {
				logutils.Error("Failed to Atoi. error: ", err)
				httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			c, err := strconv.Atoi(value.Current)
			if err != nil {
				logutils.Error("Failed to Atoi. error: ", err)
				httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response.Data.ValidDuration = d - c
		}
	}

	response.Result = true
	httpserver.Response(w, response)
}

func HandlePutLicense1(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logutils.Error("Failed to ReadAll. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := utils.DecodeSelf(string(content))
	if err != nil {
		logutils.Error("Failed to DecodeSelf. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lic := license.License{
		Auth: make(map[string]license.Auth),
	}
	err = json.Unmarshal(data, &lic.Auth)
	if err != nil {
		logutils.Error("Failed to Unmarshal. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if license.GetLicenseStatus() == base.LicenseImporting {
		httpserver.ResponseError(w, "正在导入证书", http.StatusInternalServerError)
		return
	}
	license.SetLicenseStatus(base.LicenseImporting)
	loadLicense()
	err = importLicense(&lic)
	if err != nil {
		logutils.Error("Failed to importLicense. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		loadLicense()
		checkLicense()
		return
	}
	saveLicense()
	checkLicense()
	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
}
