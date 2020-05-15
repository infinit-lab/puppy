package token

import (
	"github.com/infinit-lab/taiji/src/model/account"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/token"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"net/http"
	"time"
)

var ps *passwordSubscriber

func init() {
	logutils.Trace("Initializing controller token...")
	ps = new(passwordSubscriber)
	bus.Subscribe(base.KeyPassword, ps)
	httpserver.RegisterTokenChecker(new(tokenChecker))

	httpserver.RegisterHttpHandlerFunc(http.MethodPost, "/api/1/token", HandlePostToken1, false)
	httpserver.RegisterHttpHandlerFunc(http.MethodDelete, "/api/1/token/+", HandleDeleteToken1, false)
	go checkTokenLifetime()
}

type postToken1Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type postToken1Response struct {
	httpserver.ResponseBody
	Data string `json:"data"`
}

func duration() int {
	d := config.GetInt("token.duration")
	logutils.Trace("Get token.duration ", d)
	if d == 0 {
		d = 10 * 60
		logutils.Trace("Reset token.duration to ", d)
	}
	return d
}

func HandlePostToken1(w http.ResponseWriter, r *http.Request) {
	var request postToken1Request
	err := httpserver.GetRequestBody(r, &request)
	if err != nil {
		logutils.Error("Failed to GetRequestBody. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	isValid, err := account.IsValidAccount(request.Username, request.Password)
	if err != nil {
		logutils.Error("Failed to IsValidAccount. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isValid {
		httpserver.ResponseError(w, "用户名或密码错误", http.StatusBadRequest)
		return
	}
	t, err := token.CreateToken(request.Username, duration(), httpserver.RemoteIp(r), r)
	if err != nil {
		logutils.Error("Failed to CreateToken. error: ", err)
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response postToken1Response
	response.Result = true
	response.Data = t
	httpserver.Response(w, response)
}

func HandleDeleteToken1(w http.ResponseWriter, r *http.Request) {
	t := httpserver.GetId(r.URL.Path, "token")
	if t == "" {
		httpserver.ResponseError(w, "Url中不存在有效Token", http.StatusBadRequest)
		return
	}
	err := token.DeleteToken(t, r)
	if err != nil {
		httpserver.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response httpserver.ResponseBody
	response.Result = true
	httpserver.Response(w, response)
}

func checkTokenLifetime() {
	checkDuration := config.GetInt("token.checkDuration")
	logutils.Trace("Get token.checkDuration ", checkDuration)
	if checkDuration == 0 {
		checkDuration = 60
		logutils.Trace("Reset token.checkDuration to ", checkDuration)
	}
	for {
		time.Sleep(time.Duration(checkDuration) * time.Second)
		tokenList, err := token.GetTokenList()
		if err != nil {
			continue
		}
		now := time.Now().UTC()
		for _, t := range tokenList {
			start, err := time.Parse("2006-01-02 15:04:05", t.Time)
			if err != nil {
				_ = token.DeleteToken(t.Token, nil)
				continue
			}

			sub := now.Sub(start).Milliseconds()
			if sub < 0 || sub > (time.Duration(t.Duration)*time.Second).Milliseconds() {
				_ = token.DeleteToken(t.Token, nil)
			}
		}
	}
}
