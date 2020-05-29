package process

import (
	"archive/zip"
	"encoding/base64"
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/log"
	"github.com/infinit-lab/taiji/src/model/process"
	"github.com/infinit-lab/taiji/src/model/token"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var m *manager
var ph *processHandler
var g *guard
var s *slave

type updateManager struct {
	mutex   sync.Mutex
	updates map[int]string
}

var um *updateManager

func init() {
	logutils.Trace("Initializing controller process...")
	if config.GetBool("process.guard") {
		s = new(slave)
		s.run()
	} else {
		m = new(manager)
		m.run()

		ph = new(processHandler)
		ph.m = m
		bus.Subscribe(base.KeyProcess, ph)
		bus.Subscribe(base.KeyProcessStatus, ph)

		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process", HandleGetProcessList1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process/+", HandleGetProcess1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/process/+/operation", HandlePutProcessOperation1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process/+/status", HandleGetProcessStatusList1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process/+/status/+", HandleGetProcessStatus1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/status/+", HandleGetStatusList1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process/statistic", HandleGetProcessStatistic1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/process/+/update-file", HandlePutUpdateFile1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/update-file", HandleBatchPutUpdateFile1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process/+/config-file", HandleGetConfigFile1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodPut, "/api/1/process/+/config-file", HandlePutConfigFile1, true)
		httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/process/+/log-file", HandleGetLogFile1, true)

		g = new(guard)
		g.run()

		um = new(updateManager)
		um.updates = make(map[int]string)
	}
}

func Quit() {
	bus.Unsubscribe(base.KeyProcess, ph)
	bus.Unsubscribe(base.KeyProcessStatus, ph)
	m.quit()
}

type getProcessList1Response struct {
	httpserver.ResponseBody
	Data []*process.Process `json:"data"`
}

func HandleGetProcessList1(w http.ResponseWriter, r *http.Request) {
	var response getProcessList1Response
	var err error
	response.Data, err = process.GetProcessList()
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

func getProcessId(r *http.Request) (int, error) {
	temp := httpserver.GetId(r.URL.Path, "process")
	if temp == "" {
		return 0, errors.New("进程ID不存在")
	}
	processId, err := strconv.Atoi(temp)
	if err != nil {
		return 0, err
	}
	return processId, nil
}

type getProcess1Response struct {
	httpserver.ResponseBody
	Data *process.Process `json:"data"`
}

func HandleGetProcess1(w http.ResponseWriter, r *http.Request) {
	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response getProcess1Response
	response.Data, err = process.GetProcess(processId)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

type putProcessOperation1Request struct {
	Operation string `json:"operation"`
}

func createOperateLog(r *http.Request, processId int, operate string) {
	a := r.Header["Authorization"]
	t, err := token.GetToken(a[0])
	if err != nil {
		logutils.Error("Failed to GetToken. error: ", err)
		return
	}
	p, err := process.GetProcess(processId)
	if err != nil {
		logutils.Error("Failed to GetProcess. error: ", err)
		return
	}
	l := log.OperateLog{
		Username:    t.Username,
		Ip:          t.Ip,
		Operate:     operate,
		ProcessId:   processId,
		ProcessName: p.Name,
		Time:        time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	_ = log.CreateOperateLog(&l)
}

func HandlePutProcessOperation1(w http.ResponseWriter, r *http.Request) {
	var request putProcessOperation1Request
	if err := httpserver.GetRequestBody(r, &request); err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := m.getProcessData(processId)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	var response httpserver.ResponseBody
	logutils.Trace("Operation is ", request.Operation)
	switch request.Operation {
	case base.OperateStart:
		if p.process.Enable {
			_ = m.start(p)
		} else {
			httpserver.ResponseError(w, "该进程已禁用", http.StatusBadRequest)
			return
		}
	case base.OperateStop:
		_ = m.stop(p)
	case base.OperateRestart:
		if p.process.Enable {
			_ = m.restart(p)
		} else {
			httpserver.ResponseError(w, "该进程已禁用", http.StatusBadRequest)
			return
		}
	case base.OperateEnable:
		if !p.process.Enable {
			err = process.SetProcessEnable(p.process.Id, true, r)
			if err != nil {
				httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	case base.OperateDisable:
		if p.process.Enable {
			err = process.SetProcessEnable(p.process.Id, false, r)
			if err != nil {
				httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	default:
		httpserver.ResponseError(w, "无效操作", http.StatusBadRequest)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
	createOperateLog(r, processId, request.Operation)
}

type getProcessStatusList1Response struct {
	httpserver.ResponseBody
	Data []*process.Status `json:"data"`
}

func HandleGetProcessStatusList1(w http.ResponseWriter, r *http.Request) {
	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response getProcessStatusList1Response
	response.Data, err = process.GetStatusByProcessId(processId)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

type getProcessStatus1Response struct {
	httpserver.ResponseBody
	Data *process.Status `json:"data"`
}

func HandleGetProcessStatus1(w http.ResponseWriter, r *http.Request) {
	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	statusType := httpserver.GetId(r.URL.Path, "status")
	if statusType == "" {
		httpserver.ResponseError(w, "无效状态类型", http.StatusBadRequest)
		return
	}
	var response getProcessStatus1Response
	response.Data, err = process.GetStatus(processId, statusType)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

type getStatusList1Response struct {
	httpserver.ResponseBody
	Data []*process.Status `json:"data"`
}

func HandleGetStatusList1(w http.ResponseWriter, r *http.Request) {
	statusType := httpserver.GetId(r.URL.Path, "status")
	if statusType == "" {
		httpserver.ResponseError(w, "无效状态类型", http.StatusBadRequest)
		return
	}
	var response getStatusList1Response
	var err error
	response.Data, err = process.GetStatusByType(statusType)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = true
	httpserver.Response(w, response)
}

type getProcessStatistic1Response struct {
	httpserver.ResponseBody
	Data base.Statistic `json:"data"`
}

func HandleGetProcessStatistic1(w http.ResponseWriter, r *http.Request) {
	var response getProcessStatistic1Response
	m.mutex.Lock()
	response.Data = m.statistic
	m.mutex.Unlock()
	response.Result = true
	httpserver.Response(w, response)
}

func (m *updateManager) isUpdating(processId int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.updates[processId]
	return ok
}

func (m *updateManager) insertUpdate(processId int, updateId string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.updates[processId] = updateId
}

func (m *updateManager) eraseUpdate(processId int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.updates, processId)
}

func updateProcess(processId int, fileList []*zip.File) error {
	if um.isUpdating(processId) {
		logutils.ErrorF("Process %d is updating.", processId)
		return errors.New("正在升级")
	}
	um.insertUpdate(processId, "")
	notification := base.UpdateNotification{
		Status:  base.UpdateUpdating,
		Current: 0,
		Total:   len(fileList),
	}
	_ = bus.PublishResource(base.KeyUpdate, base.StatusUpdated, strconv.Itoa(processId), &notification, nil)

	defer func() {
		_ = bus.PublishResource(base.KeyUpdate, base.StatusUpdated, strconv.Itoa(processId), &notification, nil)
		um.eraseUpdate(processId)
	}()

	p, err := m.getProcessData(processId)
	if err != nil {
		logutils.Error("Failed to getProcessData. error: ", err)
		notification.Status = base.UpdateFail
		return errors.New("获取进程数据失败")
	}
	isStart := p.isStart
	if isStart {
		if err := m.stop(p); err != nil {
			logutils.Error("Failed to stop. error: ", err)
			notification.Status = base.UpdateFail
			return errors.New("停止进程失败")
		}
	}

	destDir := p.process.Dir
	for _, file := range fileList {
		logutils.Trace("File name is ", file.Name)
		path := filepath.Join(destDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				logutils.Error("Failed to MkdirAll. error: ", err)
				notification.Status = base.UpdateFail
				return errors.New("创建路径失败")
			}
			notification.Current++
			_ = bus.PublishResource(base.KeyUpdate, base.StatusUpdated, strconv.Itoa(processId), &notification, nil)
		} else {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				logutils.Error("Failed to MkdirAll. error: ", err)
				notification.Status = base.UpdateFail
				return errors.New("创建路径失败")
			}

			inFile, err := file.Open()
			if err != nil {
				logutils.Error("Failed to Open. error: ", err)
				notification.Status = base.UpdateFail
				return errors.New("打开文件失败")
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				logutils.Error("Failed to OpenFile. error: ", err)
				notification.Status = base.UpdateFail
				_ = inFile.Close()
				return errors.New("打开文件失败")
			}

			_, err = io.Copy(outFile, inFile)
			_ = inFile.Close()
			_ = outFile.Close()

			if err != nil {
				logutils.Error("Failed to Copy. error: ", err)
				notification.Status = base.UpdateFail
				return errors.New("拷贝文件失败")
			}

			notification.Current++
			_ = bus.PublishResource(base.KeyUpdate, base.StatusUpdated, strconv.Itoa(processId), &notification, nil)
		}
	}

	if isStart {
		if err := m.start(p); err != nil {
			logutils.Error("Failed to start. error: ", err)
		}
	}
	notification.Status = base.UpdateSuccess
	return nil
}

func HandleBatchPutUpdateFile1(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logutils.Error("Failed to ReadAll. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader := strings.NewReader(string(buffer))
	zipReader, err := zip.NewReader(reader, int64(len(buffer)))
	if err != nil {
		logutils.Error("Failed to NewReader. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	processList, err := process.GetProcessList()
	if err != nil {
		logutils.Error("Failed to GetProcessList. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go func() {
		fileMap := make(map[string][]*zip.File)
		for _, file := range zipReader.File {
			paths := strings.Split(file.Name, "/")
			for i, path := range paths {
				if path == "processes" && i+2 < len(path) && paths[i+1] != "" && paths[i+2] != "" {
					processName := paths[i+1]
					fileList, ok := fileMap[processName]
					if !ok {
						var l []*zip.File
						fileMap[processName] = l
						fileList = l
					}
					file.Name = filepath.Join(paths[i+2:]...)
					fileList = append(fileList, file)
					fileMap[processName] = fileList
					break
				}
			}
		}
		for processName, fileList := range fileMap {
			for _, p := range processList {
				if p.Name == processName {
					_ = updateProcess(p.Id, fileList)
					break
				}
			}
		}
	}()

	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
}

func HandlePutUpdateFile1(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if um.isUpdating(processId) {
		httpserver.ResponseError(w, "进程正在升级", http.StatusConflict)
		return
	}

	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logutils.Error("Failed to ReadAll. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader := strings.NewReader(string(buffer))
	zipReader, err := zip.NewReader(reader, int64(len(buffer)))
	if err != nil {
		logutils.Error("Failed to NewReader. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go func() {
		_ = updateProcess(processId, zipReader.File)
	}()

	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
	createOperateLog(r, processId, base.OperateUpdate)
}

type getConfigFile1Response struct {
	httpserver.ResponseBody
	Data string `json:"data"`
}

func HandleGetConfigFile1(w http.ResponseWriter, r *http.Request) {
	processId, err := getProcessId(r)
	if err != nil {
		logutils.Error("Failed to getProcessId. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := process.GetProcess(processId)
	if err != nil {
		logutils.Error("Failed to GetProcess. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	path := filepath.Join(p.Dir, p.ConfigFile)
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		logutils.Error("Failed to ReadAll. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response getConfigFile1Response
	response.Data = base64.StdEncoding.EncodeToString(buffer)
	response.Result = true
	httpserver.Response(w, response)
}

type putConfigFile1Request struct {
	Content string `json:"content"`
}

func HandlePutConfigFile1(w http.ResponseWriter, r *http.Request) {
	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request putConfigFile1Request
	err = httpserver.GetRequestBody(r, &request)
	if err != nil {
		logutils.Error("Failed to GetRequestBody. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	content, err := base64.StdEncoding.DecodeString(request.Content)
	if err != nil {
		logutils.Error("Failed to DecodeString. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p, err := process.GetProcess(processId)
	if err != nil {
		logutils.Error("Failed to GetProcess. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	path := filepath.Join(p.Dir, p.ConfigFile)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logutils.Error("Failed to OpenFile. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(content)
	if err != nil {
		logutils.Error("Failed to Write. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
	createOperateLog(r, processId, base.OperateConfig)
}

func HandleGetLogFile1(w http.ResponseWriter, r *http.Request) {
	processId, err := getProcessId(r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}
	p, err := m.getProcessData(processId)
	if err != nil {
		logutils.Error("Failed to getProcessData. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	file, err := os.Open(p.logFilePath)
	if err != nil {
		logutils.Error("Failed to Open. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	p.fileMutex.Lock()
	defer p.fileMutex.Unlock()

	fileStat, err := file.Stat()
	if err != nil {
		logutils.Error("Failed to Stat. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, fileName := filepath.Split(p.logFilePath)

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	w.Header().Set("File-Name", fileName)

	_, err = file.Seek(0, 0)
	if err != nil {
		logutils.Error("Failed to Seek. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = io.Copy(w, file)
	createOperateLog(r, processId, base.OperateDownloadLog)
}
