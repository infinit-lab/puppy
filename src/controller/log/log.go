package log

import (
	"github.com/infinit-lab/taiji/src/model/log"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/operate-log", HandleGetOperateLogList1, true)
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/login-log", HandleGetLoginLogList1, true)
}

type getOperateLogList1Response struct {
	httpserver.ResponseBody
	Data []*log.OperateLog `json:"data,omitempty"`
}

func HandleGetOperateLogList1(w http.ResponseWriter, r *http.Request) {
	form, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		logutils.Error("Failed to ParseQuery. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	startTime, ok := form["startTime"]
	if !ok {
		logutils.Error("Failed to Get startTime.")
		httpserver.ResponseError(w, "无开始时间", http.StatusBadRequest)
		return
	}
	stopTime, ok := form["stopTime"]
	if !ok {
		logutils.Error("Failed to Get stopTime.")
		httpserver.ResponseError(w, "无结束时间", http.StatusBadRequest)
		return
	}
	var rows int
	rowsTemp, ok := form["rows"]
	if !ok {
		logutils.Error("Failed to Get rows.")
		httpserver.ResponseError(w, "无获取行数", http.StatusBadRequest)
		return
	}
	rows, err = strconv.Atoi(rowsTemp[0])
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		httpserver.ResponseError(w, "无效获取行数", http.StatusBadRequest)
		return
	}
	var username string
	usernameTemp, ok := form["username"]
	if ok {
		username = usernameTemp[0]
	}
	var processId int
	processIdTemp, ok := form["processId"]
	if ok {
		processId, _ = strconv.Atoi(processIdTemp[0])
	}
	var offset int
	offsetTemp, ok := form["offset"]
	if ok {
		offset, _ = strconv.Atoi(offsetTemp[0])
	}

	var response getOperateLogList1Response
	response.Data, err = log.GetOperateLogList(startTime[0], stopTime[0], username, processId, rows, offset)
	if err != nil {
		logutils.Error("Failed to GetOperateLogList. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Result = true
	httpserver.Response(w, response)
}

type getLoginLogList1Response struct {
	httpserver.ResponseBody
	Data []*log.LoginLog `json:"data,omitempty"`
}

func HandleGetLoginLogList1(w http.ResponseWriter, r *http.Request) {
	form, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		logutils.Error("Failed to ParseQuery. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	startTime, ok := form["startTime"]
	if !ok {
		logutils.Error("Failed to Get startTime.")
		httpserver.ResponseError(w, "无开始时间", http.StatusBadRequest)
		return
	}
	stopTime, ok := form["stopTime"]
	if !ok {
		logutils.Error("Failed to Get stopTime.")
		httpserver.ResponseError(w, "无结束时间", http.StatusBadRequest)
		return
	}
	var rows int
	rowsTemp, ok := form["rows"]
	if !ok {
		logutils.Error("Failed to Get rows.")
		httpserver.ResponseError(w, "无获取行数", http.StatusBadRequest)
		return
	}
	rows, err = strconv.Atoi(rowsTemp[0])
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		httpserver.ResponseError(w, "无效获取行数", http.StatusBadRequest)
		return
	}
	var offset int
	offsetTemp, ok := form["offset"]
	if ok {
		offset, _ = strconv.Atoi(offsetTemp[0])
	}

	var response getLoginLogList1Response
	response.Data, err = log.GetLoginLogList(startTime[0], stopTime[0], rows, offset)
	if err != nil {
		logutils.Error("Failed to GetLoginLogList. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Result = true
	httpserver.Response(w, response)
}
