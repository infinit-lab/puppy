package token

import (
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/cache"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
	"github.com/satori/go.uuid"
	"os"
	"strings"
	"time"
)

var c *cache.Cache

func init() {
	logutils.Trace("Initializing model token...")
	table := sqlite.Table{
		Name: "token",
		Columns: []sqlite.Column{
			{
				Name:    "token",
				Type:    "CHAR(32)",
				Default: "",
				Index:   true,
				Unique:  true,
			},
			{
				Name:    "username",
				Type:    "VARCHAR(64)",
				Default: "",
				Index:   true,
			},
			{
				Name:    "ip",
				Type:    "VARCHAR(32)",
				Default: "127.0.0.1",
			},
			{
				Name:    "duration",
				Type:    "INTEGER",
				Default: "0",
			},
			{
				Name:    "time",
				Type:    "DATETIME",
				Default: "",
			},
		},
	}
	if err := base.Sqlite.InitializeTable(table); err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
	c = cache.NewCacheWithConfig()
}

type Token struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Ip       string `json:"ip"`
	Duration int    `json:"duration"`
	Time     string `json:"time"`
}

func currentDateTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func CreateToken(username string, duration int, ip string, context interface{}) (string, error) {
	token := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	_, err := base.Sqlite.Exec("INSERT INTO `token` (`token`, `username`, `ip`, `duration`, `time`) VALUES (?, ?, ?, ?, ?)",
		token, username, ip, duration, currentDateTime())
	if err != nil {
		logutils.Error("Failed to Exec. error: ", err)
		return "", err
	}
	t, _ := GetToken(token)
	_ = bus.PublishResource(base.KeyToken, base.StatusCreated, token, t, context)
	return token, nil
}

func RenewToken(token string) error {
	_, err := base.Sqlite.Exec("UPDATE `token` SET `time` = ? WHERE `token` = ?", currentDateTime(), token)
	if err != nil {
		return err
	}
	c.Erase(token)
	t, _ := GetToken(token)
	_ = bus.PublishResource(base.KeyToken, base.StatusUpdated, token, t, nil)
	return err
}

func DeleteToken(token string, context interface{}) error {
	logutils.Trace("Delete token ", token)
	t, err := GetToken(token)
	if err != nil {
		return err
	}
	_, err = base.Sqlite.Exec("DELETE FROM `token` WHERE `token` = ?", token)
	if err != nil {
		return err
	}
	c.Erase(token)
	_ = bus.PublishResource(base.KeyToken, base.StatusDeleted, token, t, context)
	return err
}

const (
	selectToken = "SELECT `token`, `username`, `ip`, `duration`, `time` FROM `token` "
)

func GetToken(token string) (*Token, error) {
	t, ok := c.Get(token)
	if ok == true {
		return t.(*Token), nil
	}
	rows, err := base.Sqlite.Query(selectToken+"WHERE `token` = ? LIMIT 1", token)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		t := new(Token)
		err := rows.Scan(&t.Token, &t.Username, &t.Ip, &t.Duration, &t.Time)
		if err != nil {
			return nil, err
		}
		t.Time = strings.Replace(t.Time, "T", " ", 1)
		t.Time = strings.Replace(t.Time, "Z", "", 1)
		c.Insert(token, t)
		return t, nil
	} else {
		return nil, errors.New("Token不存在")
	}
}

func GetTokenList() ([]*Token, error) {
	rows, err := base.Sqlite.Query(selectToken)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var tokenList []*Token
	for rows.Next() {
		t := new(Token)
		err := rows.Scan(&t.Token, &t.Username, &t.Ip, &t.Duration, &t.Time)
		if err != nil {
			return nil, err
		}
		t.Time = strings.Replace(t.Time, "T", " ", 1)
		t.Time = strings.Replace(t.Time, "Z", "", 1)
		tokenList = append(tokenList, t)
	}
	return tokenList, err
}
