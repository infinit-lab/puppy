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
		file = "./puppy.db"
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
	KeyToken int = 1
)

const (
	StatusCreated int = 1
	StatusUpdated int = 2
	StatusDeleted int = 3
)
