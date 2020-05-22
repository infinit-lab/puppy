package base

import (
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
	"os"
)

var Sqlite *sqlite.Sqlite

func init() {
	logutils.Trace("Initializing mode base...")
	file := config.GetString("sqlite.file")
	logutils.Trace("Get sqlite file is ", file)
	if file == "" {
		file = "./taiji.db"
		logutils.Warning("Reset sqlite file to ", file)
	}
	var err error
	Sqlite, err = sqlite.InitializeDatabase(file)
	if err != nil {
		logutils.Error("Failed to InitializeDatabase ", file)
		os.Exit(1)
	}
}

const (
	KeyAll        int = 0
	KeyToken      int = 1
	KeyPassword   int = 2
	KeyCpuUseRate int = 3
	KeyMemUseRate int = 4
	KeyProcess    int = 5
	//KeyProcessEnable int = 6
	KeyProcessStatus int = 7
	KeyStatistic     int = 8
	KeyUpdate        int = 9
)

const (
	StatusCreated int = 1
	StatusUpdated int = 2
	StatusDeleted int = 3
)

type MemUsage struct {
	Rate  uint32 `json:"rate"`
	Total uint64 `json:"total"`
	Avail uint64 `json:"avail"`
}

const (
	StatusTypeStarted string = "started" // 0: stopped, 1: started
)

const (
	OperateStart       string = "start"
	OperateStop        string = "stop"
	OperateRestart     string = "restart"
	OperateEnable      string = "enable"
	OperateDisable     string = "disable"
	OperateUpdate      string = "update"
	OperateConfig      string = "config"
	OperateDownloadLog string = "download_log"
)

const (
	UpdateUpdating int = 1
	UpdateSuccess  int = 2
	UpdateFail     int = 3
)

type UpdateNotification struct {
	Status  int `json:"status"`
	Current int `json:"current"`
	Total   int `json:"total"`
}

type Statistic struct {
	Total   int `json:"total"`
	Running int `json:"running"`
	Stopped int `json:"stopped"`
	Disable int `json:"disable"`
}
