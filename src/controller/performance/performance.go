package performance

import (
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"time"
)

func init() {
	go func() {
		for {
			rate, err := utils.GetCpuUseRate()
			if err != nil {
				logutils.Error("Failed to GetCpuUseRate. error: ", err)
				time.Sleep(time.Second)
				continue
			}
			_ = bus.PublishResource(base.KeyCpuUseRate, base.StatusUpdated, "", rate, nil)
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
			var usage base.MemUsage
			usage.Rate = rate
			usage.Total = total
			usage.Avail = avail
			_ = bus.PublishResource(base.KeyMemUseRate, base.StatusUpdated, "", usage, nil)
		}
	}()
}
