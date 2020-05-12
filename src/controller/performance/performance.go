package performance

import (
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"net/http"
	"time"
)

var cpuUseRate = 0
var memUseRate base.MemUsage

func init() {
	go func() {
		for {
			rate, err := utils.GetCpuUseRate()
			if err != nil {
				logutils.Error("Failed to GetCpuUseRate. error: ", err)
				time.Sleep(time.Second)
				continue
			}
			cpuUseRate = rate
			//_ = bus.PublishResource(base.KeyCpuUseRate, base.StatusUpdated, "", rate, nil)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			rate, total, avail, err := utils.GetMemoryStatus()
			if err != nil {
				logutils.Error("Failed to GetMemoryStatus. error: ", err)
				continue
			}
			memUseRate.Rate = rate
			memUseRate.Total = total
			memUseRate.Avail = avail
			//_ = bus.PublishResource(base.KeyMemUseRate, base.StatusUpdated, "", usage, nil)
		}
	}()

	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/performance/cpu", HandleGetCpu1, true)
	httpserver.RegisterHttpHandlerFunc(http.MethodGet, "/api/1/performance/mem", HandleGetMem1, true)
}

type getCpu1Response struct {
	httpserver.ResponseBody
	Data int `json:"data"`
}

func HandleGetCpu1(w http.ResponseWriter, r *http.Request) {
	var response getCpu1Response
	response.Result = true
	response.Data = cpuUseRate
	httpserver.Response(w, response)
}

type getMem1Response struct {
	httpserver.ResponseBody
	Data base.MemUsage `json:"data"`
}

func HandleGetMem1(w http.ResponseWriter, r *http.Request) {
	var response getMem1Response
	response.Result = true
	response.Data = memUseRate
	httpserver.Response(w, response)
}
