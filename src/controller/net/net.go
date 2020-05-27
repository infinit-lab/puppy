package net

import (
	n "github.com/infinit-lab/taiji/src/model/net"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"net"
	"net/http"
)

func init() {
	initNet()
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/net/interface", HandleGetNetInterfaceList1, true)
	httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/net/interface/+", HandlePutNetInterface1, true)
}

type getNetInterfaceList1Response struct {
	httpserver.ResponseBody
	Data []*utils.Adapter `json:"data"`
}

func HandleGetNetInterfaceList1(w http.ResponseWriter, r *http.Request) {
	var response getNetInterfaceList1Response
	var err error
	response.Data, err = utils.GetNetworkInfo()
	if err != nil {
		logutils.Error("Failed to GetNetworkInfo. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

func HandlePutNetInterface1(w http.ResponseWriter, r *http.Request) {
	var adapter utils.Adapter
	err := httpserver.GetRequestBody(r, &adapter)
	if err != nil {
		logutils.Error("Failed to GetRequestBody. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if net.ParseIP(adapter.Ip) == nil || net.ParseIP(adapter.Mask) == nil ||
		net.ParseIP(adapter.Gateway) == nil {
		logutils.Error("Failed to ParseIP")
		httpserver.ResponseError(w, "IP, 子网掩码或网关格式错误", http.StatusBadRequest)
		return
	}

	addr := n.Address{
		Name:    adapter.Name,
		Ip:      adapter.Ip,
		Mask:    adapter.Mask,
		Gateway: adapter.Gateway,
	}
	err = n.UpdateAddress(&addr)
	if err != nil {
		logutils.Error("Failed to UpdateAddress. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.SetAdapter(&adapter)
	if err != nil {
		logutils.Error("Failed to SetAdapter. error: ", err)
		_ = n.DeleteAddress(adapter.Name)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
}

func initNet() {
	adapters, err := utils.GetNetworkInfo()
	if err != nil {
		logutils.Error("Failed to GetNetworkInfo. error: ", err)
		return
	}
	for _, adapter := range adapters {
		addr, err := n.GetAddress(adapter.Name)
		if err != nil {
			continue
		}
		if addr.Ip != adapter.Ip || addr.Mask != adapter.Mask || addr.Gateway != adapter.Gateway {
			adapter.Ip = addr.Ip
			adapter.Mask = addr.Mask
			adapter.Gateway = addr.Gateway
			_ = utils.SetAdapter(adapter)
		}
	}
}
