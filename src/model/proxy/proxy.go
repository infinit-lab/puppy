package proxy

import (
	"database/sql"
	"errors"
	"github.com/infinit-lab/qiankun/common"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
)

func init() {
	initializeLocalServer()
	initializeRemoteHost()
}

func initializeLocalServer() {
	table := sqlite.Table {
		Name: "local_server",
		Columns: []sqlite.Column {
			{
				Name: "uuid",
				Type: "VARCHAR(64)",
				Default: "",
				Index: true,
				Unique: true,
			},
			{
				Name: "host",
				Type: "VARCHAR(256)",
				Default: "",
			},
			{
				Name: "port",
				Type: "INTEGER",
				Default: "0",
			},
			{
				Name: "description",
				Type: "VARCHAR(256)",
				Default: "",
			},
			{
				Name: "type",
				Type: "VARCHAR(256)",
				Default: "",
			},
			{
				Name: "isPublic",
				Type: "TINYINT",
				Default: "1",
			},
		},
	}
	err := base.Sqlite.InitializeTable(table)
	if err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		return
	}
}

const (
	selectLocalServer string = "SELECT `uuid`, `host`, `port`, `description`, `type`, `isPublic` FROM `local_server` "
)

func scanLocalServer(rows *sql.Rows) (*common.Server, error) {
	s := new(common.Server)
	err := rows.Scan(&s.Uuid, &s.Host, &s.Port, &s.Description, &s.Type, &s.IsPublic)
	if err != nil {
		logutils.Error("Failed to Scan. error: ", err)
		return nil, err
	}
	return s, nil
}

func GetLocalServerList() ([]*common.Server, error) {
	rows, err := base.Sqlite.Query(selectLocalServer)
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var servers []*common.Server
	for rows.Next() {
		s, err := scanLocalServer(rows)
		if err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func GetLocalServer(uuid string) (*common.Server, error) {
	rows, err := base.Sqlite.Query(selectLocalServer + " WHERE `uuid` = ? LIMIT 1", uuid)
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		logutils.Error("Can't find local server ", uuid)
		return nil, errors.New("无效本地服务UUID")
	}
	s, err := scanLocalServer(rows)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func CreateLocalServer(s *common.Server) error {
	_, err := base.Sqlite.Exec("INSERT INTO `local_server` (`uuid`, `host`, `port`, `description`, `type`, `isPublic`) " +
		"VALUES (?, ?, ?, ?, ?, ?)", s.Uuid, s.Host, s.Port, s.Description, s.Type, s.IsPublic)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

func UpdateLocalServer(s *common.Server) error {
	_, err := base.Sqlite.Exec("UPDATE `local_server` SET `host` = ?, `port` = ?, `description` = ? " +
		"WHERE `uuid` = ?", s.Host, s.Port, s.Description, s.Uuid)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

func DeleteLocalServer(uuid string) error {
	_, err := base.Sqlite.Exec("DELETE FROM `local_server` WHERE `uuid` = ?", uuid)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

func initializeRemoteHost() {
	table := sqlite.Table {
		Name: "remote_host",
		Columns: []sqlite.Column {
			{
				Name: "address",
				Type: "VARCHAR(256)",
				Default: "",
				Index: true,
				Unique: true,
			},
			{
				Name: "description",
				Type: "VARCHAR(512)",
				Default: "",
			},
		},
	}
	err := base.Sqlite.InitializeTable(table)
	if err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		return
	}
}

type RemoteHost struct {
	Address string `json:"address"`
	Description string `json:"description"`
}

func GetRemoteHostList() ([]*RemoteHost, error) {
	rows, err := base.Sqlite.Query("SELECT `address`, `description` FROM `remote_host`")
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var remoteHostList []*RemoteHost
	for rows.Next() {
		r := new(RemoteHost)
		err := rows.Scan(&r.Address, &r.Description)
		if err != nil {
			logutils.Error("Failed to Scan. error: ", err)
			return nil, err
		}
		remoteHostList = append(remoteHostList, r)
	}
	return remoteHostList, nil
}

func CreateRemoteHost(r *RemoteHost) error {
	_, err := base.Sqlite.Exec("INSERT INTO `remote_host` (`address`, `description`) VALUES (?, ?)",
		r.Address, r.Description)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

func DeleteRemoteHost(addr string) error {
	_, err := base.Sqlite.Exec("DELETE FROM `remote_host` WHERE `address` = ?", addr)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

