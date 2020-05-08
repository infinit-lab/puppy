package process

import (
	"encoding/json"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"testing"
	"time"
)

type processHandler struct {
	Id int
}

func (h *processHandler) Handle(key int, value *bus.Resource) {
	logutils.TraceF("Key is %d, status is %d", key, value.Status)
	if key == base.KeyProcess || key == base.KeyProcessEnable {
		switch value.Status {
		case base.StatusCreated, base.StatusUpdated, base.StatusDeleted:
			p, ok := value.Data.(*Process)
			if !ok {
				logutils.Error("Failed to convert data to process.")
				return
			}
			data, _ := json.Marshal(p)
			logutils.Trace("Process is %s", string(data))
			h.Id = p.Id
		default:
			return
		}
	}

	if key == base.KeyProcessStatus {
		s, ok := value.Data.(*Status)
		if !ok {
			logutils.Error("Failed to convert data to status.")
			return
		}
		data, _ := json.Marshal(s)
		logutils.Trace("Status is ", string(data))
	}
}

var ph *processHandler

func TestInit(t *testing.T) {
	os.Args = append(os.Args, "log.level=3")
	config.Exec()
	ph = new(processHandler)
	bus.Subscribe(base.KeyProcess, ph)
	bus.Subscribe(base.KeyProcessEnable, ph)
	bus.Subscribe(base.KeyProcessStatus, ph)
}

func compareProcess(t *testing.T, process1 *Process, process2 *Process) {
	if process2.Id != ph.Id || process2.Name != process1.Name || process2.Path != process1.Path ||
		process2.Dir != process1.Dir || process2.Config != process1.Config || process2.Enable != process1.Enable {
		p1, _ := json.Marshal(process1)
		p2, _ := json.Marshal(process2)
		t.Errorf("Expect is %s, actual is %s", string(p1), string(p2))
	}
}

func TestCreateProcess(t *testing.T) {
	process := new(Process)
	process.Name = "CreateName"
	process.Path = "CreatePath"
	process.Dir = "CreateDir"
	process.Config = "CreateConfig"
	process.Enable = true
	err := CreateProcess(process, nil)
	if err != nil {
		t.Error("Failed to CreateProcess. error: ", err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	if ph.Id == 0 {
		t.Error("Process id should not be 0.")
		return
	}
	p, err := GetProcess(ph.Id)
	if err != nil {
		t.Error("Failed to GetProcess. error: ", err)
		return
	}
	compareProcess(t, process, p)
}

func TestUpdateProcess(t *testing.T) {
	process := new(Process)
	process.Id = ph.Id
	process.Name = "UpdateName"
	process.Path = "UpdatePath"
	process.Dir = "UpdateDir"
	process.Config = "UpdateConfig"
	process.Enable = false

	if err := UpdateProcess(process.Id, process, nil); err != nil {
		t.Error("Failed to UpdateProcess. error: ", err)
		return
	}
	p, err := GetProcess(process.Id)
	if err != nil {
		t.Error("Failed to GetProcess. error: ", err)
		return
	}
	compareProcess(t, process, p)
}

func TestGetProcessList(t *testing.T) {
	processList, err := GetProcessList()
	if err != nil {
		t.Error("Failed to GetProcessList. error: ", err)
	}
	for _, process := range processList {
		data, _ := json.Marshal(process)
		logutils.Trace("Get process ", string(data))
	}
}

func TestSetProcessEnable(t *testing.T) {
	if err := SetProcessEnable(ph.Id, true, nil); err != nil {
		t.Error("Failed to SetProcessEnable. error: ", err)
		return
	}
	process, err := GetProcess(ph.Id)
	if err != nil {
		t.Error("Failed to GetProcess. error: ", err)
		return
	}
	if process.Enable != true {
		t.Error("Failed to SetProcessEnable.")
	}
}

func TestDeleteProcess(t *testing.T) {
	err := DeleteProcess(ph.Id, nil)
	if err != nil {
		t.Error("Failed to DeleteProcess. error: ", err)
		return
	}
	_, err = GetProcess(ph.Id)
	if err == nil {
		t.Error("Should not get process ", ph.Id)
	}
	time.Sleep(100 * time.Millisecond)
}

func TestUpdateStatus(t *testing.T) {
	status := new(Status)
	status.ProcessId = 1
	status.Type = "Type"
	status.Value = "Create"
	if err := UpdateStatus(status, nil); err != nil {
		t.Error("Failed to UpdateStatus. error: ", err)
		return
	}
	status2, err := GetStatus(status.ProcessId, status.Type)
	if err != nil {
		t.Error("Failed to GetStatus. error: ", err)
		return
	}
	if status.ProcessId != status2.ProcessId || status.Type != status2.Type || status.Value != status2.Value {
		s1, _ := json.Marshal(status)
		s2, _ := json.Marshal(status2)
		t.Errorf("Expect is %s, actual is %s", string(s1), string(s2))
		return
	}
	status.Value = "Update"
	if err := UpdateStatus(status, nil); err != nil {
		t.Error("Failed to UpdateStatus. error: ", err)
		return
	}
	status2, err = GetStatus(status.ProcessId, status.Type)
	if err != nil {
		t.Error("Failed to GetStatus. error: ", err)
		return
	}
	if status.ProcessId != status2.ProcessId || status.Type != status2.Type || status.Value != status2.Value {
		s1, _ := json.Marshal(status)
		s2, _ := json.Marshal(status2)
		t.Errorf("Expect is %s, actual is %s", string(s1), string(s2))
		return
	}
}

func TestGetStatusList(t *testing.T) {
	statusList, err := GetStatusList()
	if err != nil {
		t.Error("Failed to GetStatusList. error: ", err)
		return
	}
	for _, status := range statusList {
		data, _ := json.Marshal(status)
		logutils.Trace("Status is ", string(data))
	}
	statusList, err = GetStatusByProcessId(1)
	if err != nil {
		t.Error("Failed to GetStatusList. error: ", err)
		return
	}
	for _, status := range statusList {
		data, _ := json.Marshal(status)
		logutils.Trace("Status is ", string(data))
	}
	statusList, err = GetStatusByType("Type")
	if err != nil {
		t.Error("Failed to GetStatusList. error: ", err)
		return
	}
	for _, status := range statusList {
		data, _ := json.Marshal(status)
		logutils.Trace("Status is ", string(data))
	}
}

func TestDeleteStatus(t *testing.T) {
	err := DeleteStatus(1, "Type", nil)
	if err != nil {
		t.Error("Failed to DeleteStatus. error: ", err)
		return
	}
	_, err = GetStatus(1, "Type")
	if err == nil {
		t.Error("Should not get status.")
		return
	}
	time.Sleep(100 * time.Millisecond)
}
