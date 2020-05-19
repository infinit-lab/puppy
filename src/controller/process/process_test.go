package process

import (
	"encoding/json"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type testProcessHandler struct {
}

func (h *testProcessHandler) Handle(key int, value *bus.Resource) {
	data, err := json.Marshal(value.Data)
	if err != nil {
		logutils.Error("Failed to Marshal. error: ", err)
		return
	}
	logutils.TraceF("Key is %d, data is %s", key, string(data))
}

var tPH *testProcessHandler

func TestInit(t *testing.T) {
	tPH = new(testProcessHandler)
	bus.Subscribe(base.KeyProcess, tPH)
	bus.Subscribe(base.KeyProcessStatus, tPH)
	bus.Subscribe(base.KeyStatistic, tPH)

	go func() {
		_ = httpserver.ListenAndServe()
	}()
}

var processId int

func TestHandleGetProcessList1(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	code, rsp, _, err := utils.Request(http.MethodGet, "http://127.0.0.1:8088/api/1/process", nil, nil, "")
	if err != nil {
		t.Error("Failed to Request. error: ", err)
		return
	}
	if code != 200 {
		t.Error("Failed to Request. code: ", code)
		return
	}
	var response getProcessList1Response
	err = json.Unmarshal(rsp, &response)
	if err != nil {
		t.Error("Failed to Unmarshal. error: ")
		return
	}
	if !response.Result {
		t.Error("Failed to Request. result: ", response.Result)
		return
	}
	for i, r := range response.Data {
		processId = r.Id
		data, err := json.Marshal(r)
		if err != nil {
			continue
		}
		logutils.TraceF("%d. Process is %s", i+1, string(data))
	}
	time.Sleep(1000 * time.Millisecond)
}

func putProcessOperation1(operation string, t *testing.T) {
	request := putProcessOperation1Request{
		Operation: operation,
	}
	code, rsp, _, err := utils.Request(http.MethodPut, "http://127.0.0.1:8088/api/1/process/"+strconv.Itoa(processId)+"/operation", nil, request, "")
	if err != nil {
		t.Error("Failed to Request. error: ", err)
		return
	}
	if code != http.StatusOK {
		t.Error("Failed to Request. code: ", code)
		return
	}
	var response httpserver.ResponseBody
	err = json.Unmarshal(rsp, &response)
	if err != nil {
		t.Error("Failed to Unmarshal. error: ", err)
		return
	}
	if response.Result != true {
		t.Error("Failed to Request. result: ", response.Result)
		return
	}
}

func TestHandlePutProcessOperation1(t *testing.T) {
	putProcessOperation1(base.OperateRestart, t)
	time.Sleep(time.Second)
	putProcessOperation1(base.OperateStop, t)
	time.Sleep(time.Second)
	putProcessOperation1(base.OperateStart, t)
	time.Sleep(time.Second)
	putProcessOperation1(base.OperateDisable, t)
	time.Sleep(time.Second)
	putProcessOperation1(base.OperateEnable, t)
	time.Sleep(time.Second)
	putProcessOperation1(base.OperateStop, t)
	time.Sleep(time.Second)
}

func TestHandleGetProcessStatusList1(t *testing.T) {
	code, rsp, _, err := utils.Request(http.MethodGet, "http://127.0.0.1:8088/api/1/process/"+strconv.Itoa(processId)+"/status", nil, nil, "")
	if err != nil {
		t.Error("Failed to Request. error: ", err)
		return
	}
	if code != 200 {
		t.Error("Failed to Request. code: ", code)
		return
	}
	var response getProcessStatusList1Response
	err = json.Unmarshal(rsp, &response)
	if err != nil {
		t.Error("Failed to Unmarshal. error: ", err)
		return
	}
	if response.Result != true {
		t.Error("Failed to Request. result: ", response.Result)
		return
	}
	for i, r := range response.Data {
		data, err := json.Marshal(r)
		if err != nil {
			continue
		}
		logutils.TraceF("%d. Status is %s", i+1, string(data))
	}
}

func TestHandleGetProcessStatus1(t *testing.T) {
	code, rsp, _, err := utils.Request(http.MethodGet, "http://127.0.0.1:8088/api/1/process/"+strconv.Itoa(processId)+"/status/started", nil, nil, "")
	if err != nil {
		t.Error("Failed to Request. error: ", err)
		return
	}
	if code != 200 {
		t.Error("Failed to Request. code: ", code)
		return
	}
	var response getProcessStatus1Response
	err = json.Unmarshal(rsp, &response)
	if err != nil {
		t.Error("Failed to Unmarshal. error: ", err)
		return
	}
	if response.Result != true {
		t.Error("Failed to Request. result: ", response.Result)
		return
	}
	data, err := json.Marshal(response.Data)
	if err != nil {
		t.Error("Failed to Marshal. error: ", err)
		return
	}
	logutils.Trace("Status is ", string(data))
}

func TestHandleGetStatusList1(t *testing.T) {
	code, rsp, _, err := utils.Request(http.MethodGet, "http://127.0.0.1:8088/api/1/status/started", nil, nil, "")
	if err != nil {
		t.Error("Failed to Request. error: ", err)
		return
	}
	if code != 200 {
		t.Error("Failed to Request. code: ", code)
		return
	}
	var response getStatusList1Response
	err = json.Unmarshal(rsp, &response)
	if err != nil {
		t.Error("Failed to Unmarshal. error: ", err)
		return
	}
	if response.Result != true {
		t.Error("Failed to Request. result: ", response.Result)
		return
	}
	for i, r := range response.Data {
		data, err := json.Marshal(r)
		if err != nil {
			continue
		}
		logutils.TraceF("%d. Status is %s", i+1, string(data))
	}
}

func TestHandleGetProcessStatistic1(t *testing.T) {
	code, rsp, _, err := utils.Request(http.MethodGet, "http://127.0.0.1:8088/api/1/process/statistic", nil, nil, "")
	if err != nil {
		t.Error("Failed to Request. error: ", err)
		return
	}
	if code != 200 {
		t.Error("Failed to Request. code: ", code)
		return
	}
	var response getProcessStatistic1Response
	err = json.Unmarshal(rsp, &response)
	if err != nil {
		t.Error("Failed to Unmarshal. error: ", err)
		return
	}
	if response.Result != true {
		t.Error("Failed to Request. result: ", response.Result)
		return
	}
	data, err := json.Marshal(response.Data)
	if err != nil {
		t.Error("Failed to Marshal. error: ", err)
		return
	}
	logutils.Trace("Statistic is ", string(data))
}
