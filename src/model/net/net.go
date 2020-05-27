package net

import (
	"database/sql"
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
	"os"
	"strconv"
)

func init() {
	table := sqlite.Table{
		Name: "net",
		Columns: []sqlite.Column{
			{
				Name:    "name",
				Type:    "VARCHAR(256)",
				Default: "",
				Index:   true,
				Unique:  true,
			},
			{
				Name:    "ip",
				Type:    "VARCHAR(32)",
				Default: "",
			},
			{
				Name:    "mask",
				Type:    "VARCHAR(32)",
				Default: "",
			},
			{
				Name:    "gateway",
				Type:    "VARCHAR(32)",
				Default: "",
			},
		},
	}
	if err := base.Sqlite.InitializeTable(table); err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
}

type Address struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Ip      string `json:"ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
}

const selectSql = "SELECT `id`, `name`, `ip`, `mask`, `gateway` FROM `net` "

func scanAddress(rows *sql.Rows) (*Address, error) {
	a := new(Address)
	if err := rows.Scan(&a.Id, &a.Name, &a.Ip, &a.Mask, &a.Gateway); err != nil {
		logutils.Error("Failed to Scan. error: ", err)
		return nil, err
	}
	return a, nil

}

func GetAddress(name string) (*Address, error) {
	rows, err := base.Sqlite.Query(selectSql+" WHERE `name` = ? LIMIT 1", name)
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		return scanAddress(rows)
	} else {
		return nil, errors.New("无效地址")
	}
}

func GetAddressList() ([]*Address, error) {
	rows, err := base.Sqlite.Query(selectSql)
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var addresses []*Address
	for rows.Next() {
		a, err := scanAddress(rows)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func UpdateAddress(a *Address) error {
	_, err := GetAddress(a.Name)
	var sqlString string
	if err != nil {
		sqlString = "INSERT INTO `net` (`name`, `ip`, `mask`, `gateway`) VALUES (?, ?, ?, ?)"
	} else {
		sqlString = "UPDATE `net` SET `name` = ?, `ip` = ?, `mask` = ?, `gateway` = ? WHERE `id` = " +
			strconv.Itoa(a.Id)
	}
	_, err = base.Sqlite.Exec(sqlString, a.Name, a.Ip, a.Mask, a.Gateway)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}

func DeleteAddress(name string) error {
	_, err := base.Sqlite.Exec("DELETE FROM `net` WHERE `name` = ?", name)
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return err
	}
	return nil
}
