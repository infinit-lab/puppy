package log

import (
	"database/sql"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
	"os"
	"strconv"
	"strings"
)

func init() {
	initOperateLogTable()
	initLoginLogTable()
}

type OperateLog struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Ip          string `json:"ip"`
	Operate     string `json:"operate"`
	ProcessId   int    `json:"processId"`
	ProcessName string `json:"processName"`
	Time        string `json:"time"`
}

func initOperateLogTable() {
	table := sqlite.Table{
		Name: "operate_log",
		Columns: []sqlite.Column{
			{
				Name:    "username",
				Type:    "VARCHAR(64)",
				Default: "",
				Index:   true,
			},
			{
				Name:    "ip",
				Type:    "VARCHAR(32)",
				Default: "",
			},
			{
				Name:    "operate",
				Type:    "VARCHAR(32)",
				Default: "",
				Index:   true,
			},
			{
				Name:    "processId",
				Type:    "INTEGER",
				Default: "0",
				Index:   true,
			},
			{
				Name:    "processName",
				Type:    "VARCHAR(64)",
				Default: "",
			},
			{
				Name:    "time",
				Type:    "DATETIME",
				Default: "",
				Index:   true,
			},
		},
	}
	err := base.Sqlite.InitializeTable(table)
	if err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
}

func CreateOperateLog(l *OperateLog) error {
	_, err := base.Sqlite.Exec("INSERT INTO `operate_log` (`username`, `ip`, "+
		"`operate`, `processId`, `processName`, `time`) VALUES (?, ?, ?, ?, ?, ?)",
		l.Username, l.Ip, l.Operate, l.ProcessId, l.ProcessName, l.Time)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

const (
	selectOperateLog = "SELECT `id`, `username`, `ip`, `operate`, `processId`, `processName`, `time` " +
		"FROM `operate_log` "
)

func scanOperateLog(rows *sql.Rows) (*OperateLog, error) {
	l := new(OperateLog)
	err := rows.Scan(&l.Id, &l.Username, &l.Ip, &l.Operate, &l.ProcessId, &l.ProcessName, &l.Time)
	if err != nil {
		logutils.Error("Failed to Scan. error: ", err)
		return nil, err
	}
	l.Time = strings.Replace(l.Time, "T", " ", 1)
	l.Time = strings.Replace(l.Time, "Z", "", 1)
	return l, nil
}

func GetOperateLogList(startTime, stopTime string, username string, processId int, rows, offset int) ([]*OperateLog, error) {
	sqlString := selectOperateLog + " WHERE `time` BETWEEN '" + startTime + "' AND '" + stopTime + "'"
	if username != "" {
		sqlString += " AND `username` = " + username
	}
	if processId != 0 {
		sqlString += " AND `processId` = " + strconv.Itoa(processId)
	}
	sqlString += " ORDER BY `id` DESC"
	sqlString += " LIMIT " + strconv.Itoa(rows)
	sqlString += " OFFSET " + strconv.Itoa(offset)

	r, err := base.Sqlite.Query(sqlString)
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()
	var logList []*OperateLog
	for r.Next() {
		l, err := scanOperateLog(r)
		if err != nil {
			return nil, err
		}
		logList = append(logList, l)
	}
	return logList, nil
}

type LoginLog struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Ip       string `json:"ip"`
	IsLogin  bool   `json:"isLogin"`
	Time     string `json:"time"`
}

func initLoginLogTable() {
	table := sqlite.Table{
		Name: "login_log",
		Columns: []sqlite.Column{
			{
				Name:    "username",
				Type:    "VARCHAR(64)",
				Default: "",
				Index:   true,
			},
			{
				Name:    "ip",
				Type:    "VARCHAR(32)",
				Default: "",
				Index:   true,
			},
			{
				Name:    "isLogin",
				Type:    "TINYINT",
				Default: "0",
				Index:   true,
			},
			{
				Name:    "time",
				Type:    "DATETIME",
				Default: "",
				Index:   true,
			},
		},
	}
	if err := base.Sqlite.InitializeTable(table); err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
}

func CreateLoginLog(l *LoginLog) error {
	_, err := base.Sqlite.Exec("INSERT INTO `login_log` (`username`, `ip`, `isLogin`, `time`) "+
		"VALUES (?, ?, ?, ?)", l.Username, l.Ip, l.IsLogin, l.Time)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

const (
	selectLoginLog = "SELECT `id`, `username`, `ip`, `isLogin`, `time` FROM `login_log` "
)

func scanLoginLog(r *sql.Rows) (*LoginLog, error) {
	l := new(LoginLog)
	if err := r.Scan(&l.Id, &l.Username, &l.Ip, &l.IsLogin, &l.Time); err != nil {
		logutils.Error("Failed to Scan. error: ", err)
		return nil, err
	}
	l.Time = strings.Replace(l.Time, "T", " ", 1)
	l.Time = strings.Replace(l.Time, "Z", "", 1)
	return l, nil
}

func GetLoginLogList(startTime, stopTime string, rows, offset int) ([]*LoginLog, error) {
	sqlString := selectLoginLog + " WHERE `time` BETWEEN '" + startTime + "' AND '" + stopTime + "' " +
		" ORDER BY `id` DESC LIMIT " + strconv.Itoa(rows) + " OFFSET " + strconv.Itoa(offset)
	r, err := base.Sqlite.Query(sqlString)
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()
	var loginLogList []*LoginLog
	for r.Next() {
		l, err := scanLoginLog(r)
		if err != nil {
			return nil, err
		}
		loginLogList = append(loginLogList, l)
	}
	return loginLogList, nil
}
