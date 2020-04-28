package account
import (
	"errors"
	"github.com/infinit-lab/puppy/src/model/base"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/sqlite"
	"os"
)

func init() {
	logutils.Trace("Initializing model account...")
	table := sqlite.Table {
		Name: "account",
		Columns: []sqlite.Column {
			{
				Name: "username",
				Type: "VARCHAR(64)",
				Default: "",
				Index: true,
				Unique: true,
			},
			{
				Name: "password",
				Type: "CHAR(32)",
				Default: "",
			},
		},
	}
	if err := base.Sqlite.InitializeTable(table); err != nil {
		logutils.Error("Failed to InitializeTable. error: ", err)
		os.Exit(1)
	}
	_ = initializeAccount()
}

func initializeAccount() error {
	rows, err := base.Sqlite.Query("SELECT `username` FROM `account` WHERE `username` = 'admin'")
	if err != nil {
		logutils.Error("Failed to Query. error: ", err)
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		_, err := base.Sqlite.Exec("INSERT INTO `account` (`username`, `password`) VALUES ('admin', '21232f297a57a5a743894a0e4a801fc3')")
		if err != nil {
			logutils.Error("Failed to Exec. error: ", err)
			return nil
		}
	}
	return nil
}

func IsValidAccount(username string, password string) (bool, error) {
	rows, err := base.Sqlite.Query("SELECT `username` FROM `account` WHERE `username` = ? AND `password` = ?", username, password)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		return false, nil
	}
	return true, nil
}

func ChangePassword(username string, originPassword string, newPassword string) error {
	ret, err := base.Sqlite.Exec("UPDATE `account` SET `password` = ? WHERE `username` = ? AND `password` = ?", newPassword, username, originPassword)
	if err == nil {
		rows, err := ret.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("用户名或密码错误")
		}
	}
	return err
}

