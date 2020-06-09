package net

import (
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/log"
	n "github.com/infinit-lab/taiji/src/model/net"
	"github.com/infinit-lab/taiji/src/model/token"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"net"
	"net/http"
	"time"
)

func init() {
	logutils.Trace("Initializing controller net...")
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

func SetAdapter(adapter *utils.Adapter) error {
	if net.ParseIP(adapter.Ip) == nil || net.ParseIP(adapter.Mask) == nil ||
		net.ParseIP(adapter.Gateway) == nil {
		logutils.Error("Failed to ParseIP")
		return errors.New("IP, 子网掩码或网关格式错误")
	}

	addr := n.Address{
		Name:    adapter.Name,
		Ip:      adapter.Ip,
		Mask:    adapter.Mask,
		Gateway: adapter.Gateway,
	}
	err := n.UpdateAddress(&addr)
	if err != nil {
		logutils.Error("Failed to UpdateAddress. error: ", err)
		return err
	}
	err = utils.SetAdapter(adapter)
	if err != nil {
		logutils.Error("Failed to SetAdapter. error: ", err)
		_ = n.DeleteAddress(adapter.Name)
		return err
	}
	return nil
}

func HandlePutNetInterface1(w http.ResponseWriter, r *http.Request) {
	var adapter utils.Adapter
	err := httpserver.GetRequestBody(r, &adapter)
	if err != nil {
		logutils.Error("Failed to GetRequestBody. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = SetAdapter(&adapter)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)

	a := r.Header["Authorization"]
	t, err := token.GetToken(a[0])
	if err != nil {
		logutils.Error("Failed to GetToken. error: ", err)
		return
	}
	l := log.OperateLog{
		Username:    t.Username,
		Ip:          t.Ip,
		Operate:     base.OperateConfigNet,
		ProcessId:   0,
		ProcessName: adapter.Name,
		Time:        time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	_ = log.CreateOperateLog(&l)
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
