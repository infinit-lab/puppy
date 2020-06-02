package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/license"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	self, err := utils.GetMachineFingerprint()
	if err != nil {
		logutils.Error("Failed to GetMachineFingerprint. error: ", err)
		return
	}
	if self != "__L5Kxm5bmDcRMCxDjlHIQ==" {
		logutils.Error("Invalid machine fingerprint.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please input fingerprint> ")
	buffer, _, err := reader.ReadLine()
	if err != nil {
		logutils.Error("Failed to ReadLine. error: ", err)
		return
	}
	fingerprint := string(buffer)

	fmt.Print("Please input is forever(y/n)> ")
	buffer, _, err = reader.ReadLine()
	if err != nil {
		logutils.Error("Failed to ReadLine. error: ", err)
		return
	}
	isForever := false
	if string(buffer) == "y" {
		isForever = true
		generateLicense(fingerprint, isForever, nil)
		return
	}

	fmt.Print("Please input year> ")
	buffer, _, err = reader.ReadLine()
	if err != nil {
		logutils.Error("Failed to ReadLine. error: ", err)
		return
	}
	year, err := strconv.Atoi(string(buffer))
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return
	}

	fmt.Print("Please input month> ")
	buffer, _, err = reader.ReadLine()
	if err != nil {
		logutils.Error("Failed to ReadLine. error: ", err)
		return
	}
	month, err := strconv.Atoi(string(buffer))
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return
	}

	fmt.Print("Please input day> ")
	buffer, _, err = reader.ReadLine()
	if err != nil {
		logutils.Error("Failed to ReadLine. error: ", err)
		return
	}
	day, err := strconv.Atoi(string(buffer))
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 23:59:59", year, month, day))
	if err != nil {
		logutils.Error("Failed to Parse. error: ", err)
	}
	generateLicense(fingerprint, isForever, &t)
}

func generateLicense(fingerprint string, isForever bool, datetime *time.Time) {
	auths := make(map[string]license.Auth)
	var a license.Auth

	u := uuid.NewV4().String()
	a.Type = base.AuthUuid
	a.ValueType = base.ValueTypeString
	a.Value = []string{u}
	auths[a.Type] = a

	a.Type = base.AuthForever
	a.ValueType = base.ValueTypeBool
	a.Value = []string{strconv.FormatBool(isForever)}
	auths[a.Type] = a

	if isForever == false {
		now := time.Now().UTC()
		duration := datetime.Sub(now).Seconds()
		if duration <= 0 {
			logutils.Error("Invalid date.")
			return
		}
		a.Type = base.AuthDatetime
		a.ValueType = base.ValueTypeDatetime
		a.Value = []string{datetime.Format("2006-01-02 15:04:05")}
		auths[a.Type] = a

		a.Type = base.AuthDuration
		a.ValueType = base.ValueTypeInt
		a.Value = []string{strconv.Itoa(int(duration))}
		auths[a.Type] = a
	}

	data, err := json.Marshal(auths)
	if err != nil {
		logutils.Error("Failed to Marshal. error: ", err)
		return
	}
	content, err := utils.Encode(fingerprint, data)
	if err != nil {
		logutils.Error("Failed to Encode. error: ", err)
		return
	}
	filename := fingerprint + "_" + u + ".txt"
	err = ioutil.WriteFile(filename, []byte(content), os.ModePerm)
	if err != nil {
		logutils.Error("Failed to WriteFile. error: ", err)
		return
	}
	path, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Generate license ", filename)
	} else {
		fmt.Println("Generate license ", path)
	}
}
