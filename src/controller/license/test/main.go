package main

import (
	"encoding/json"
	_ "github.com/infinit-lab/taiji/src/controller/license"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/license"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	_ = writeLicense()
}

func writeLicense() error {
	var l license.License
	l.Auth = make(map[string]license.Auth)
	var a license.Auth

	a.Type = base.AuthUuid
	a.ValueType = base.ValueTypeString
	u := uuid.NewV4().String()
	a.Value = []string{u}
	l.Auth[a.Type] = a

	a.Type = base.AuthForever
	a.Name = "永久有效"
	a.ValueType = base.ValueTypeBool
	a.Value = []string{"false"}
	l.Auth[a.Type] = a

	a.Type = base.AuthDatetime
	a.Name = "有效日期"
	a.ValueType = base.ValueTypeDatetime
	a.Value = []string{
		time.Now().UTC().Add(600 * time.Second).Format("2006-01-02 15:04:05"),
	}
	l.Auth[a.Type] = a

	a.Type = base.AuthDuration
	a.Name = "有效时间"
	a.ValueType = base.ValueTypeInt
	a.Value = []string{"600"}
	l.Auth[a.Type] = a

	data, _ := json.Marshal(l.Auth)
	content, err := utils.EncodeSelf(data)
	if err != nil {
		logutils.Error("Failed to EncodeSelf. error: ", err)
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logutils.Error("Failed to Abs. error: ", err)
		return err
	}
	fingerprint, err := utils.GetMachineFingerprint()
	if err != nil {
		logutils.Error("Failed to GetMachineFingerprint. error: ", err)
		return err
	}
	filename := fingerprint + "_" + u + ".txt"
	path := filepath.Join(dir, filename)
	err = ioutil.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		logutils.Error("Failed to WriteFile. error: ", err)
		return err
	}
	return nil
}
