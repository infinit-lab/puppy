package process

import (
	"database/sql"
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/cache"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
	"os"
	"strconv"
)

var processCache *cache.Cache
var statusCache *cache.Cache

func init() {
	logutils.Trace("Initializing model process...")
	processTable := sqlite.Table{
		Name: "process",
		Columns: []sqlite.Column{
			{
				Name:    "name",
				Type:    "VARCHAR(64)",
				Default: "",
			},
			{
				Name:    "path",
				Type:    "VARCHAR(256)",
				Default: "",
			},
			{
				Name:    "dir",
				Type:    "VARCHAR(256)",
				Default: "",
			},
			{
				Name:    "config",
				Type:    "TEXT",
				Default: "",
			},
			{
				Name:    "enable",
				Type:    "TINYINT",
				Default: "1",
			},
			{
				Name:    "pid",
				Type:    "INTEGER",
				Default: "0",
			},
			{
				Name: "startTime",
				Type: "VARCHAR(32)",
				Default: "",
			},
			{
				Name: "configFile",
				Type: "VARCHAR(256)",
				Default: "",
			},
		},
	}
	if err := base.Sqlite.InitializeTable(processTable); err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
	processCache = cache.NewCacheWithConfig()

	processStatusTable := sqlite.Table{
		Name: "process_status",
		Columns: []sqlite.Column{
			{
				Name:    "processId",
				Type:    "INTEGER",
				Default: "0",
				Index:   true,
			},
			{
				Name:    "type",
				Type:    "VARCHAR(256)",
				Default: "",
				Index:   true,
			},
			{
				Name:    "value",
				Type:    "VARCHAR(256)",
				Default: "",
			},
		},
	}
	if err := base.Sqlite.InitializeTable(processStatusTable); err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
	statusCache = cache.NewCacheWithConfig()
}

type Process struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Dir    string `json:"path"`
	Config string `json:"config"`
	Enable bool   `json:"enable"`
	Pid    int    `json:"pid"`
	StartTime string `json:"startTime"`
	ConfigFile string `json:"configFile"`
}

type Status struct {
	ProcessId int    `json:"processId"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

const (
	selectProcess = "SELECT `id`, `name`, `path`, `dir`, `config`, `enable`, `pid`, `startTime`, `configFile` " +
		"FROM `process` "
)

func scanProcess(rows *sql.Rows) (*Process, error) {
	p := new(Process)
	err := rows.Scan(&p.Id, &p.Name, &p.Path, &p.Dir, &p.Config, &p.Enable, &p.Pid, &p.StartTime, &p.ConfigFile)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetProcess(id int) (*Process, error) {
	p, ok := processCache.Get(strconv.Itoa(id))
	if ok == true {
		return p.(*Process), nil
	}
	rows, err := base.Sqlite.Query(selectProcess+"WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		p, err := scanProcess(rows)
		if err != nil {
			return nil, err
		}
		processCache.Insert(strconv.Itoa(id), p)
		return p, nil
	} else {
		return nil, errors.New("进程不存在")
	}
}

func GetProcessList() ([]*Process, error) {
	pList, ok := processCache.Get("list")
	if ok {
		return pList.([]*Process), nil
	}
	rows, err := base.Sqlite.Query(selectProcess)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var processList []*Process
	for rows.Next() {
		p, err := scanProcess(rows)
		if err != nil {
			return nil, err
		}
		processList = append(processList, p)
	}
	processCache.Insert("list", processList)
	return processList, nil
}

func CreateProcess(p *Process, context interface{}) error {
	ret, err := base.Sqlite.Exec("INSERT INTO `process` (`name`, `path`, `dir`, `config`, `enable`, `pid`, `configFile`) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?)", p.Name, p.Path, p.Dir, p.Config, p.Enable, p.Pid, p.ConfigFile)
	var id int64
	if err != nil {
		return err
	} else {
		rows, err := ret.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("创建进程失败")
		}
		id, err = ret.LastInsertId()
		if err != nil {
			return err
		}
	}
	processCache.Erase("list")
	process, _ := GetProcess(int(id))
	_ = bus.PublishResource(base.KeyProcess, base.StatusCreated, strconv.Itoa(int(id)), process, context)
	return nil
}

func UpdateProcess(id int, p *Process, context interface{}) error {
	_, err := base.Sqlite.Exec("UPDATE `process` "+
		"SET `name` = ?, `path` = ?, `dir` = ?, `config` = ?, `enable` = ?, `pid` = ?, `startTime` = ?, `configFile` = ? " +
		"WHERE `id` = ?",
		p.Name, p.Path, p.Dir, p.Config, p.Enable, p.Pid, p.StartTime, p.ConfigFile, id)
	if err != nil {
		return nil
	}
	processCache.Erase("list")
	processCache.Erase(strconv.Itoa(id))
	process, _ := GetProcess(id)
	_ = bus.PublishResource(base.KeyProcess, base.StatusUpdated, strconv.Itoa(id), process, context)
	return nil
}

func DeleteProcess(id int, context interface{}) error {
	process, err := GetProcess(id)
	if err != nil {
		return err
	}
	_, err = base.Sqlite.Exec("DELETE FROM `process` WHERE `id` = ?", id)
	if err != nil {
		return nil
	}
	processCache.Erase("list")
	processCache.Erase(strconv.Itoa(id))
	_ = bus.PublishResource(base.KeyProcess, base.StatusDeleted, strconv.Itoa(id), process, context)
	return nil
}

func SetProcessEnable(id int, enable bool, context interface{}) error {
	_, err := base.Sqlite.Exec("UPDATE `process` SET `enable` = ? WHERE `id` = ?", enable, id)
	if err != nil {
		return err
	}
	processCache.Erase("list")
	processCache.Erase(strconv.Itoa(id))
	process, _ := GetProcess(id)
	_ = bus.PublishResource(base.KeyProcess, base.StatusUpdated, strconv.Itoa(id), process, context)
	return nil
}

const (
	selectStatus = "SELECT `processId`, `type`, `value` FROM `process_status` "
)

func scanStatus(rows *sql.Rows) (*Status, error) {
	s := new(Status)
	if err := rows.Scan(&s.ProcessId, &s.Type, &s.Value); err != nil {
		logutils.Error("Failed to Scan. error: ", err)
		return nil, err
	}
	return s, nil
}

func statusKey(processId int, statusType string) string {
	return strconv.Itoa(processId) + "_" + statusType
}

func GetStatus(processId int, statusType string) (*Status, error) {
	key := statusKey(processId, statusType)
	s, ok := statusCache.Get(key)
	if ok {
		return s.(*Status), nil
	}
	rows, err := base.Sqlite.Query(selectStatus+" WHERE `processId` = ? AND `type` = ? LIMIT 1",
		processId, statusType)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		s, err := scanStatus(rows)
		if err != nil {
			return nil, err
		}
		statusCache.Insert(key, s)
		return s, nil
	} else {
		return nil, errors.New("状态不存在")
	}
}

func GetStatusByProcessId(processId int) ([]*Status, error) {
	rows, err := base.Sqlite.Query(selectStatus+" WHERE `processId` = ?", processId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var statusList []*Status
	for rows.Next() {
		s, err := scanStatus(rows)
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, s)
	}
	return statusList, nil
}

func GetStatusByType(statusType string) ([]*Status, error) {
	rows, err := base.Sqlite.Query(selectStatus+" WHERE `type` = ?", statusType)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var statusList []*Status
	for rows.Next() {
		s, err := scanStatus(rows)
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, s)
	}
	return statusList, nil
}

func GetStatusList() ([]*Status, error) {
	rows, err := base.Sqlite.Query(selectStatus)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var statusList []*Status
	for rows.Next() {
		s, err := scanStatus(rows)
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, s)
	}
	return statusList, nil
}

func UpdateStatus(status *Status, context interface{}) error {
	s, err := GetStatus(status.ProcessId, status.Type)
	if err != nil {
		_, err = base.Sqlite.Exec("INSERT INTO `process_status` (`processId`, `type`, `value`) VALUES (?, ?, ?)",
			status.ProcessId, status.Type, status.Value)
		if err != nil {
			return err
		}
	} else if s.Value == status.Value {
		return nil
	} else {
		_, err = base.Sqlite.Exec("UPDATE `process_status` SET `value` = ? WHERE `processId` = ? AND `type` = ?",
			status.Value, status.ProcessId, status.Type)
		if err != nil {
			return err
		}
	}
	statusCache.Erase(statusKey(status.ProcessId, status.Type))
	s, _ = GetStatus(status.ProcessId, status.Type)
	_ = bus.PublishResource(base.KeyProcessStatus, base.StatusUpdated, strconv.Itoa(status.ProcessId), s, context)
	return nil
}

func DeleteStatus(processId int, statusType string, context interface{}) error {
	s, err := GetStatus(processId, statusType)
	if err != nil {
		return err
	}
	_, err = base.Sqlite.Exec("DELETE FROM `process_status` WHERE `processId` = ? AND `type` = ?",
		processId, statusType)
	if err != nil {
		return err
	}
	statusCache.Erase(statusKey(processId, statusType))
	_ = bus.PublishResource(base.KeyProcessStatus, base.StatusDeleted, strconv.Itoa(processId), s, context)
	return nil
}
