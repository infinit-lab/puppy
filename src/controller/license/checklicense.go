package license

import (
	"encoding/json"
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/license"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func FilePath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logutils.Error("Failed to Abs. error: ", err)
		return "", err
	}
	fingerprint, err := utils.GetMachineFingerprint()
	if err != nil {
		logutils.Error("Failed to GetMachineFingerprint. error: ", err)
		return "", err
	}
	filename := "." + fingerprint + ".txt"
	path := filepath.Join(dir, filename)
	return path, nil
}

func loadLicense() {
	path, err := FilePath()
	if err != nil {
		return
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logutils.Error("Failed to ReadFile. error: ", err)
		return
	}
	data, err := utils.DecodeSelf(string(content))
	if err != nil {
		logutils.Error("Failed to DecodeSelf. error: ", err)
		return
	}
	logutils.Trace(string(data))

	l.Auth = make(map[string]license.Auth)
	err = json.Unmarshal(data, &l.Auth)
	if err != nil {
		logutils.Error("Failed to Unmarshal. error: ", err)
		return
	}
}

func checkLicense() {
	status := base.LicenseUnauthorized
	defer license.SetLicenseStatus(status)

	isForever, err := checkForever()
	if err != nil {
		return
	}
	if isForever == true {
		status = base.LicenseAuthorized
		return
	}

	isInDatetime, err := checkInDatetime()
	if err != nil {
		return
	}
	if isInDatetime == false {
		return
	}

	isInDuration, err := checkDuration()
	if err != nil {
		return
	}
	if isInDuration == true {
		status = base.LicenseAuthorized
	}
}

func checkForever() (bool, error){
	forever, ok := l.Auth[base.AuthForever]
	if !ok {
		logutils.Error("Failed to find ", base.AuthForever)
		return false, errors.New("未找到是否永久使用")
	}
	isForever, err := strconv.ParseBool(forever.Current)
	if err != nil {
		logutils.Error("Failed to ParseBool. error: ", err)
		return false, err
	}
	return isForever, nil
}

func checkInDatetime() (bool, error) {
	datetime, ok := l.Auth[base.AuthDatetime]
	if !ok {
		logutils.Error("Failed to find ", base.AuthDatetime)
		return false, errors.New("未找到有效日期")
	}
	if len(datetime.Value) == 0 {
		logutils.Error("Datetime value is 0.")
		return false, errors.New("未找到有效日期数据")
	}
	deadline, err := time.Parse("2006-01-02 15:04:05", datetime.Value[0])
	if err != nil {
		logutils.Error("Failed to Parse. error: ", err)
		return false, err
	}
	currentTime, err := time.Parse("2006-01-02 15:04:05", datetime.Current)
	if err != nil {
		logutils.Error("Failed to Parse. error: ", err)
		return false, err
	}
	subTime := deadline.Sub(currentTime)
	if subTime <= 0 {
		logutils.Error("End of trial!!!")
		return false, nil
	}
	return true, nil
}

func checkDuration() (bool, error) {
	duration, ok := l.Auth[base.AuthDuration]
	if !ok {
		logutils.Error("Failed to find ", base.AuthDuration)
		return false, errors.New("未找到有效时间")
	}
	if len(duration.Value) == 0 {
		logutils.Error("Duration value is 0.")
		return false, errors.New("未找到有效时间数据")
	}
	d, err := strconv.Atoi(duration.Value[0])
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return false, err
	}
	c, err := strconv.Atoi(duration.Current)
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return false, err
	}
	if d - c <= 0 {
		logutils.Error("End of trial!!!")
		return false, nil
	}
	return true, nil
}

func updateLicenseLoop() {
	for {
		time.Sleep(time.Minute)
		if license.GetLicenseStatus() != base.LicenseAuthorized {
			continue
		}
		loadLicense()
		updateLicense()
		saveLicense()
		checkLicense()
	}
}

func saveLicense() {
	data, err := json.Marshal(l.Auth)
	if err != nil {
		logutils.Error("Failed to Marshal. error: ", err)
		return
	}
	content, err := utils.EncodeSelf(data)
	if err != nil {
		logutils.Error("Failed to EncodeSelf. error: ", err)
		return
	}
	path, err := FilePath()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		logutils.Error("Failed to WriteFile. error: ", err)
		return
	}
}

func updateLicense() {
	isForever, err := checkForever()
	if err != nil {
		return
	}
	if isForever == true {
		return
	}
	delta1, err := getDatetimeDelta()
	if err != nil {
		return
	}
	delta2, err := getDurationDelta()
	if err != nil {
		return
	}
	var delta int
	if delta1 < delta2 {
		delta = delta1
	} else {
		delta = delta2
	}
	if delta < 0 {
		delta = 0
	}
	updateDatetime(delta)
	updateDuration(delta)
}

func getDatetimeDelta() (int, error) {
	datetime, ok := l.Auth[base.AuthDatetime]
	if !ok {
		logutils.Error("Failed to find ", base.AuthDatetime)
		return 0, errors.New("未找到有效日期")
	}
	if len(datetime.Value) == 0 {
		logutils.Error("Datetime value is 0.")
		return 0, errors.New("未找到有效日期数据")
	}
	deadline, err := time.Parse("2006-01-02 15:04:05", datetime.Value[0])
	if err != nil {
		logutils.Error("Failed to Parse. error: ", err)
		return 0, err
	}
	currentTime, err := time.Parse("2006-01-02 15:04:05", datetime.Current)
	if err != nil {
		logutils.Error("Failed to Parse. error: ", err)
		return 0, err
	}
	currentTime = currentTime.Add(time.Minute)
	return int(deadline.Sub(currentTime)), nil
}

func updateDatetime(delta int) {
	datetime, ok := l.Auth[base.AuthDatetime]
	if !ok {
		logutils.Error("Failed to find ", base.AuthDatetime)
		return
	}
	if len(datetime.Value) == 0 {
		logutils.Error("Datetime value is 0.")
		return
	}
	now := time.Now().Local()
	datetime.Current = now.Format("2006-01-02 15:04:05")
	datetime.Value[0] = now.Add(time.Duration(delta)).Local().Format("2006-01-02 15:04:05")
	l.Auth[base.AuthDatetime] = datetime
}

func getDurationDelta() (int, error) {
	duration, ok := l.Auth[base.AuthDuration]
	if !ok {
		logutils.Error("Failed to find ", base.AuthDuration)
		return 0, errors.New("未找到有效时间")
	}
	if len(duration.Value) == 0 {
		logutils.Error("Duration value is 0.")
		return 0, errors.New("未找到有效时间数据")
	}
	d, err := strconv.Atoi(duration.Value[0])
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return 0, err
	}
	c, err := strconv.Atoi(duration.Current)
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return 0, err
	}
	return d - (c + 60), nil
}

func updateDuration(delta int) {
	duration, ok := l.Auth[base.AuthDuration]
	if !ok {
		logutils.Error("Failed to find ", base.AuthDuration)
		return
	}
	if len(duration.Value) == 0 {
		logutils.Error("Duration value is 0.")
		return
	}
	d, err := strconv.Atoi(duration.Value[0])
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return
	}
	duration.Current = strconv.Itoa(d - delta)
	l.Auth[base.AuthDuration] = duration
}

func importLicense(lic *license.License) error {
	temp := license.License {
		Auth: make(map[string]license.Auth),
	}
	var err error
	temp.Auth[base.AuthUuid], err = compareUuid(lic)
	if err != nil {
		return err
	}

	foreverAuth, ok := lic.Auth[base.AuthForever]
	if !ok {
		logutils.Error("Failed to Find ", base.AuthForever)
		return errors.New("未找到是否永久有效")
	}
	if len(foreverAuth.Value) == 0 {
		logutils.Error("foreverAuth.Value is 0")
		return errors.New("未找到是否永久有效数据")
	}
	foreverAuth.Current = foreverAuth.Value[0]
	temp.Auth[base.AuthForever] = foreverAuth

	isForever, err := strconv.ParseBool(foreverAuth.Current)
	if err != nil {
		logutils.Error("Failed to ParseBool. error: ", err)
		return err
	}
	if isForever == false {
		datetime, duration, err := calculateDuration(lic)
		if err != nil {
			return err
		}
		temp.Auth[base.AuthDatetime] = datetime
		temp.Auth[base.AuthDuration] = duration
	}

	for key, value := range lic.Auth {
		if key == base.AuthUuid || key == base.AuthForever ||
			key == base.AuthDatetime || key == base.AuthDuration {
			continue
		}
		tempAuth, ok := l.Auth[key]
		if ok {
			value.Current = tempAuth.Current
		}
		temp.Auth[key] = tempAuth
	}
	l.Auth = temp.Auth
	return nil
}

func compareUuid(lic *license.License) (license.Auth, error) {
	uuidAuth1, ok := l.Auth[base.AuthUuid]
	if !ok {
		logutils.Error("Failed to Find ", base.AuthUuid)
		return license.Auth{}, errors.New("未找到证书UUID")
	}
	if len(uuidAuth1.Value) == 0 {
		logutils.Error("uuidAuth1.Value is 0")
		return license.Auth{}, errors.New("未找到证书UUID数据")
	}
	uuidAuth2, ok := lic.Auth[base.AuthUuid]
	if !ok {
		logutils.Error("Failed to Find ", base.AuthUuid)
		return license.Auth{}, errors.New("未找到证书UUID")
	}
	if len(uuidAuth2.Value) == 0 {
		logutils.Error("uuidAuth2.Value is 0")
		return license.Auth{}, errors.New("为找到证书UUID数据")
	}
	for _, u := range uuidAuth1.Value {
		if u == uuidAuth2.Value[0] {
			return license.Auth{}, errors.New("重复导入证书")
		}
	}
	uuidAuth1.Value = append(uuidAuth1.Value, uuidAuth2.Value[0])
	return uuidAuth1, nil
}

func calculateDuration(lic *license.License) (datetime license.Auth, duration license.Auth, err error) {
	datetimeAuth, ok := lic.Auth[base.AuthDatetime]
	if !ok {
		logutils.Error("Failed to Find ", base.AuthDatetime)
		err = errors.New("未找到有效日期")
		return
	}
	if len(datetimeAuth.Value) == 0 {
		logutils.Error("datetimeAuth.Value is 0")
		err = errors.New("未找到有效日期数据")
		return
	}
	durationAuth, ok := lic.Auth[base.AuthDuration]
	if !ok {
		logutils.Error("Failed to Find ", base.AuthDuration)
		err = errors.New("未找到有效时间")
		return
	}
	if len(durationAuth.Value) == 0 {
		logutils.Error("durationAuth.Value is 0")
		err = errors.New("未找到有效时间数据")
		return
	}

	now := time.Now().Local()
	dt, err := time.Parse("2006-01-02 15:04:05", datetimeAuth.Value[0])
	if err != nil{
		logutils.Error("Failed to Parse. error: ", err)
		return
	}
	dur, err := strconv.Atoi(durationAuth.Value[0])
	if err != nil {
		logutils.Error("Failed to Atoi. error: ", err)
		return
	}
	delta := int(dt.Sub(now))
	if delta > dur {
		delta = dur
	}
	if delta < 0 {
		delta = 0
	}
	datetime = datetimeAuth
	datetime.Current = now.Format("2006-01-02 15:04:05")
	datetime.Value[0] = now.Add(time.Duration(delta)).Local().Format("2006-01-02 15:04:05")

	duration = durationAuth
	duration.Current = "0"
	duration.Value[0] = strconv.Itoa(delta)
	return
}
