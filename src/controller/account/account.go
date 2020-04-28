package account

import (
	"github.com/infinit-lab/puppy/src/model/account"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"net/http"
)

func init() {
	logutils.Trace("Initializing controller account...")
	httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/password/+", HandlePutPassword1, true)
}

type putPassword1Request struct {
	Origin string `json:"origin"`
	New    string `json:"new"`
}

func HandlePutPassword1(w http.ResponseWriter, r *http.Request) {
	var request putPassword1Request
	err := httpserver.GetRequestBody(r, &request)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := httpserver.GetId(r.URL.Path, "password")
	err = account.ChangePassword(username, request.Origin, request.New, r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
}
