package system

import (
	"github.com/infinit-lab/yolanda/httpserver"
	"net/http"
)

var Version string
var CommitId string
var BuildTime string

func init() {
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/version", HandleGetVersion1, false)
}

type GetVersion1Response struct {
	httpserver.ResponseBody
	Version string `json:"version"`
	CommitId string `json:"commitId"`
	BuildTime string `json:"buildTime"`
}

func HandleGetVersion1(w http.ResponseWriter, r *http.Request) {
	var response GetVersion1Response
	response.Result = true
	response.Version = Version
	response.CommitId = CommitId
	response.BuildTime = BuildTime
	httpserver.Response(w, response)
}
