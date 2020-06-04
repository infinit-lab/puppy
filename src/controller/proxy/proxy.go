package proxy

import (
	"github.com/infinit-lab/qiankun/common"
	"github.com/infinit-lab/qiankun/kun"
	"github.com/infinit-lab/taiji/src/model/proxy"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"net/http"
	"strings"
)

func init() {
	servers, err := proxy.GetLocalServerList()
	if err == nil {
		for _, server := range servers {
			kun.AddServer(*server)
		}
	}

	hosts := strings.Split(config.GetString("kun.hosts"), ",")
	for _, host := range hosts {
		if host == "" {
			continue
		}
		_ = kun.AddQian(host)
	}

	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/proxy/local-server", HandleGetLocalServerList1, false)
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/proxy/local-server/+", HandleGetLocalServer1, false)
	httpserver.RegisterHttpHandlerFunc(http.MethodPost, "/api/1/proxy/local-server", HandleCreateLocalServer1, false)
	httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/proxy/local-server/+", HandleUpdateLocalServer1, false)
	httpserver.RegisterHttpHandlerFunc(http.MethodDelete, "/api/1/proxy/local-server/+", HandleDeleteLocalServer1, false)
}

type getLocalServerList1Response struct {
	httpserver.ResponseBody
	Data []*common.Server `json:"data"`
}

func HandleGetLocalServerList1(w http.ResponseWriter, r *http.Request) {
	var response getLocalServerList1Response
	var err error
	response.Data, err = proxy.GetLocalServerList()
	if err != nil {
		logutils.Error("Failed to GetLocalServerList. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

func getLocalServerUuid(w http.ResponseWriter, r *http.Request) string {
	uuid := httpserver.GetId(r.URL.Path, "local-server")
	if uuid == "" {
		httpserver.ResponseError(w, "无效本地服务UUID", http.StatusBadRequest)
	}
	return uuid
}

type getLocalServer1Response struct {
	httpserver.ResponseBody
	Data *common.Server `json:"data"`
}

func HandleGetLocalServer1(w http.ResponseWriter, r *http.Request) {
	uuid := getLocalServerUuid(w, r)
	if uuid == "" {
		return
	}
	var response getLocalServer1Response
	var err error
	response.Data, err = proxy.GetLocalServer(uuid)
	if err != nil {
		logutils.Error("Failed to GetLocalServer. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

func HandleCreateLocalServer1(w http.ResponseWriter, r *http.Request) {
	s := new(common.Server)
	err := httpserver.GetRequestBody(r, s)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = proxy.CreateLocalServer(s)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kun.AddServer(*s)
	response := httpserver.ResponseBody {
		Result: true,
	}
	httpserver.Response(w, response)
}

func HandleUpdateLocalServer1(w http.ResponseWriter, r *http.Request) {
	uuid := getLocalServerUuid(w, r)
	if uuid == "" {
		return
	}
	s := new(common.Server)
	err := httpserver.GetRequestBody(r, s)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = proxy.UpdateLocalServer(s)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kun.UpdateServer(*s)
	response := httpserver.ResponseBody {
		Result: true,
	}
	httpserver.Response(w, response)
}

func HandleDeleteLocalServer1(w http.ResponseWriter, r *http.Request) {
	uuid := getLocalServerUuid(w, r)
	if uuid == "" {
		return
	}
	err := proxy.DeleteLocalServer(uuid)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kun.DeleteServer(uuid)
	response := httpserver.ResponseBody {
		Result: true,
	}
	httpserver.Response(w, response)
}
