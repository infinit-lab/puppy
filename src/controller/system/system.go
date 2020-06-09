package system

import (
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"net/http"
)

type Version struct {
	Version string `json:"version"`
	CommitId string `json:"commitId"`
	BuildTime string `json:"buildTime"`
}

var version Version

func init() {
	logutils.Trace("Initializing controller system...")
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/version", HandleGetVersion1, false)
}

type GetVersion1Response struct {
	httpserver.ResponseBody
	Data struct {
		Version string `json:"version"`
		CommitId string `json:"commitId"`
		BuildTime string `json:"buildTime"`
	} `json:"data"`
}

func HandleGetVersion1(w http.ResponseWriter, r *http.Request) {
	var response GetVersion1Response
	response.Result = true
	response.Data.Version = version.Version
	response.Data.CommitId = version.CommitId
	response.Data.BuildTime = version.BuildTime
	httpserver.Response(w, response)
}

func SetVersion(v *Version) {
	version = *v
}

func GetVersion() Version {
	return version
}
